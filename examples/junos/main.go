package main

import (
	"fmt"
	"log"

	"github.com/metajar/netconf/netconf"
	"github.com/metajar/netconf/netconf/message"
)

func main() {
	// Initialize the netconf client to connect to the JUNOS device.
	client, err := netconf.NewClient("172.31.255.5:22", "admin", "Password")
	if err != nil {
		log.Fatalf("Failed to create netconf client: %v", err)
	}
	defer client.Close()

	// Get interface information. Example using raw for sending raw RPC calls for Gets.
	interfaceFilter := `<get-interface-information><detail/></get-interface-information>`
	interfaceInfo, err := client.Raw(interfaceFilter)
	if err != nil {
		log.Fatalf("Failed to get interface information: %v", err)
	}
	fmt.Println("Interface Information:")
	fmt.Println(interfaceInfo)

	// Define the configuration filter and the configuration to be set
	configFilter := `<configuration><interfaces><interface><name>et-0/0/0</name></interface></interfaces></configuration>`
	newConfig := `<configuration>
		<interfaces>
			<interface>
				<name>et-0/0/0</name>
				<description>This is the new description</description>
			</interface>
		</interfaces>
	</configuration>`

	// Retrieve the current configuration
	currentConfig, err := client.GetConfig(message.DatastoreRunning, "subtree", configFilter)
	if err != nil {
		log.Fatalf("Failed to get current configuration: %v", err)
	}
	fmt.Println("Current Configuration:")
	fmt.Println(currentConfig)

	// Lock the configuration
	_, err = client.Lock()
	if err != nil {
		log.Fatalf("Failed to lock configuration: %v", err)
	}
	defer func() {
		if _, err := client.UnLock(); err != nil {
			log.Fatalf("Failed to unlock configuration: %v", err)
		}
	}()

	// Edit the configuration
	_, err = client.Edit(newConfig)
	if err != nil {
		log.Fatalf("Failed to edit configuration: %v", err)
	}

	// Commit the configuration
	_, err = client.Commit()
	if err != nil {
		log.Fatalf("Failed to commit configuration: %v", err)
	}

	// Retrieve the updated configuration
	updatedConfig, err := client.GetConfig(message.DatastoreRunning, "subtree", configFilter)
	if err != nil {
		log.Fatalf("Failed to get updated configuration: %v", err)
	}
	fmt.Println("Updated Configuration:")
	fmt.Println(updatedConfig)
}
