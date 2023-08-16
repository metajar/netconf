package main

import (
	"fmt"
	"github.com/metajar/netconf/netconf"
	"log"
)

func main() {
	// Connect to the Arista device
	newArista, err := netconf.NewClient("192.168.88.9:830", "grpc", "53cret", netconf.ARISTATYPE)
	if err != nil {
		log.Fatalln(err)
	}
	defer newArista.Close()

	// Get the running configuration
	running, err := newArista.GetRunning()
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(running)

	// Setting the payload of just simply changing the
	// hostname.
	payload := `
<system xmlns="http://openconfig.net/yang/system">
	<config>
		<hostname>Arista101</hostname> 
	</config>
</system>`

	// Locking the candidate config for configuration
	// operation.
	lock, err := newArista.Lock()
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("Locked:", lock)
	defer func() {
		_, err := newArista.UnLock()
		if err != nil {
			log.Fatalln(err)
		}
	}()

	// Send the Edit to the device with the configuration
	// payload.
	edit, err := newArista.Edit(payload)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("Edit", edit)

	// Commit the configuration after we are complete.
	commit, err := newArista.Commit()
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("Commit:", commit)
}
