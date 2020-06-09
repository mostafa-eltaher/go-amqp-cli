package main

import (
	"bufio"
	"context"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/Azure/go-amqp"
	"gopkg.in/yaml.v2"
)

// YamlConfig is exported.
type YamlConfig struct {
	Spec struct {
		Connection struct {
			Container      string `yaml:"container"`
			Authentication struct {
				Type     string `yaml:"type"`
				Username string `yaml:"username"`
				Password string `yaml:"password"`
			} `yaml:"authentication"`
		} `yaml:"connection"`
		Sessions []struct {
			Links []struct {
				Role          string `yaml:"role"`
				Source        string `yaml:"source"`
				Target        string `yaml:"target"`
				InitialCredit uint32 `yaml:"initialCredit"`
			} `yaml:"links"`
		} `yaml:"sessions"`
	} `yaml:"spec"`
}

func main() {
	fmt.Println("Parsing YAML file")
	reader := bufio.NewReader(os.Stdin)

	yamlConfig, err := parseConf()

	if err != nil {
		fmt.Printf("Error parsing YAML file: %s\n", err)
		panic(err)
	}
	container := yamlConfig.Spec.Connection.Container
	username := yamlConfig.Spec.Connection.Authentication.Username
	password := yamlConfig.Spec.Connection.Authentication.Password

	log.Printf("Dialing AMQP server: %s", container)
	client, err := amqp.Dial(container,
		amqp.ConnSASLPlain(username, password),
	)
	if err != nil {
		log.Fatal("Dialing AMQP server:", err)
	}
	defer client.Close()
	// Open a session
	for _, sessionConfig := range yamlConfig.Spec.Sessions {
		session, err := client.NewSession()
		if err != nil {
			log.Fatal("Creating AMQP session:", err)
		}

		ctx := context.Background()
		defer session.Close(ctx)
		for _, linkConfig := range sessionConfig.Links {

			if linkConfig.Role == "sender" {
				// Send a message
				{
					// Create a sender
					sender, err := session.NewSender(
						amqp.LinkTargetAddress(linkConfig.Target),
					)
					if err != nil {
						log.Fatal("Creating sender link:", err)
					}

					//ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
					//ctx, cancel := context.WithTimeout(ctx, 5*time.Second)

					// Send message
					err = sender.Send(ctx, amqp.NewMessage([]byte("Hello!")))
					if err != nil {
						log.Fatal("Sending message:", err)
					}

					func() {
						for {
							fmt.Print("-> ")
							text, _ := reader.ReadString('\n')
							fmt.Printf("-> SENDING ...%s", text)
							err = sender.Send(ctx, amqp.NewMessage([]byte(text)))
							if err != nil {
								log.Fatal("Sending message:", err)
							}
						}
					}()
					//defer sender.Close(ctx)
					//defer cancel()
				}
			} else if linkConfig.Role == "receiver" {
				// Continuously read messages
				{
					// Create a receiver
					receiver, err := session.NewReceiver(
						amqp.LinkSourceAddress(linkConfig.Source),
						amqp.LinkCredit(linkConfig.InitialCredit),
					)
					if err != nil {
						log.Fatal("Creating receiver link:", err)
					}
					defer func() {
						ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
						receiver.Close(ctx)
						cancel()
					}()
					func() {
						for {
							// Receive next message
							msg, err := receiver.Receive(ctx)
							if err != nil {
								log.Fatal("Reading message from AMQP:", err)
							}

							// Accept message
							msg.Accept()

							fmt.Printf("Message received: %s\n", msg.GetData())
						}
					}()
				}
			}
		}
	}

	fmt.Printf("Result: %v\n", yamlConfig)
}

func parseConf() (*YamlConfig, error) {
	var fileName string
	flag.StringVar(&fileName, "f", "", "YAML file to parse.")
	flag.Parse()

	if fileName == "" {
		fmt.Println("Please provide yaml file by using -f option")

	}

	yamlFile, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Printf("Error reading YAML file: %s\n", err)
		return nil, errors.New("Error reading YAML file")
	}

	var yamlConfig YamlConfig
	err = yaml.Unmarshal(yamlFile, &yamlConfig)

	return &yamlConfig, err
}
