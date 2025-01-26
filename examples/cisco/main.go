package main

import (
	"fmt"
	"log"

	"github.com/metajar/netconf/netconf"
)

func main() {
	// We setup a new netconf client with a CISCOTYPE to connect to the cisco device.
	c, err := netconf.NewClient("172.20.20.2:830", "clab", "clab@123")
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
  <interfaces xmlns="http://openconfig.net/yang/interfaces">
    <interface>
      <name>GigabitEthernet0/0/0/0</name>
      <config>
        <name>GigabitEthernet0/0/0/0</name>
        <description>Hello There Friend</description>
      </config>
      <subinterfaces>
        <subinterface>
          <index>0</index>
          <ipv4 xmlns="http://openconfig.net/yang/interfaces/ip">
            <addresses>
              <address>
                <ip>88.88.88.1</ip>
                <config>
                  <ip>88.88.88.1</ip>
                  <prefix-length>30</prefix-length>
                </config>
              </address>
            </addresses>
          </ipv4>
        </subinterface>
      </subinterfaces>
    </interface>
  </interfaces>
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
}
