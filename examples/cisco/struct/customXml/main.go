package main

import (
	"encoding/xml"
	"fmt"
)

type InterfaceConfigurations struct {
	XMLName xml.Name                 `xml:"interface-configurations"`
	Configs []InterfaceConfiguration `xml:"interface-configuration"`
}

type InterfaceConfiguration struct {
	Active        string   `xml:"active"`
	Description   string   `xml:"description"`
	InterfaceName string   `xml:"interface-name"`
	Shutdown      Shutdown `xml:"shutdown"`
}

type Shutdown struct {
	Op string `xml:"operation,attr"`
}

func (s Shutdown) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	start.Attr = []xml.Attr{
		{Name: xml.Name{Local: "xmlns:nc"}, Value: "urn:ietf:params:xml:ns:netconf:base:1.0"},
		{Name: xml.Name{Local: "nc:operation"}, Value: s.Op},
	}
	e.EncodeToken(start)
	e.EncodeToken(start.End())
	return nil
}

func main() {
	config := InterfaceConfigurations{
		Configs: []InterfaceConfiguration{
			{
				Active:        "act",
				Description:   "This Interface",
				InterfaceName: "GigabitEthernet0/0/0/0",
				Shutdown: Shutdown{
					Op: "delete",
				},
			},
			{
				Active:        "act",
				Description:   "That Interface",
				InterfaceName: "GigabitEthernet0/0/0/1",
				Shutdown: Shutdown{
					Op: "delete",
				},
			},
		},
	}

	output, err := xml.MarshalIndent(config, "", "    ")
	if err != nil {
		panic(err)
	}
	fmt.Println(xml.Header + string(output))
}
