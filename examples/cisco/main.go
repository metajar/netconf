package main

import (
	"fmt"
	"log"

	"github.com/metajar/netconf/netconf"
)

func main() {
	// We setup a new netconf client with a CISCOTYPE to connect to the cisco device.
	c, err := netconf.NewClient("172.20.20.2:830", "grpc", "53cret", netconf.CISCOTYPE)
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()
	running, err := c.GetRunning()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(running)

	payload := `
	<interface-configurations xmlns="http://cisco.com/ns/yang/Cisco-IOS-XR-ifmgr-cfg">
        <interface-configuration>
            <active>act</active>
			<description>This is A description</description>
            <interface-name>GigabitEthernet0/0/0/0</interface-name>
		<shutdown xmlns:nc="urn:ietf:params:xml:ns:netconf:base:1.0" nc:operation="delete"/>
        </interface-configuration>
    </interface-configurations>
`
	_, err = c.Lock()
	if err != nil {
		log.Fatalln(err)
	}

	defer func() {
		_, err := c.UnLock()
		if err != nil {
			log.Fatalln(err)
		}
	}()

	_, err = c.Edit(payload)
	if err != nil {
		log.Fatalln(err)
	}

	_, err = c.Commit()
	if err != nil {
		log.Fatalln(err)
	}

	platformFilter := `
<components xmlns="http://openconfig.net/yang/platform">
    <component>
        <state>
        </state>
    </component>
</components>
`

	response, err := c.Get(platformFilter)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(response)
}
