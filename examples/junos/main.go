package main

import (
	"fmt"
	"log"

	"github.com/metajar/netconf/netconf"
)

func main() {
	// We setup a new netconf client with a JUNOS to connect to the junos device.
	c, err := netconf.NewClient("172.31.255.5:22", "admin", "Password", netconf.JUNOS)
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()
	filter := `<configuration>
      <interfaces/>
    </configuration>`

	resp, err := c.GetConfig("candidate", "subtree", filter)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("%s\n", resp)
}
