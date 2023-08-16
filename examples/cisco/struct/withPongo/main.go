package main

import (
	"fmt"
	"github.com/flosch/pongo2/v4"
)

type InterfaceConfiguration struct {
	Active        string
	Description   string
	InterfaceName string
	ShutdownOp    string
}

func main() {
	templateStr := `
<interface-configurations xmlns="http://cisco.com/ns/yang/Cisco-IOS-XR-ifmgr-cfg">
    {% for interface in interfaces -%}
    <interface-configuration>
        <active>{{ interface.Active }}</active>
        <description>{{ interface.Description }}</description>
        <interface-name>{{ interface.InterfaceName }}</interface-name>
        <shutdown xmlns:nc="urn:ietf:params:xml:ns:netconf:base:1.0" nc:operation="{{ interface.ShutdownOp }}"></shutdown>
    </interface-configuration>
    {% endfor -%}
</interface-configurations>
`

	// Populate the InterfaceConfiguration structs
	interfaces := []InterfaceConfiguration{
		{
			Active:        "act",
			Description:   "Fthis",
			InterfaceName: "GigabitEthernet0/0/0/0",
			ShutdownOp:    "delete",
		},
		{
			Active:        "act2",
			Description:   "Fthat",
			InterfaceName: "GigabitEthernet0/0/0/1",
			ShutdownOp:    "none",
		},
	}

	ctx := pongo2.Context{
		"interfaces": interfaces,
	}

	tmpl, err := pongo2.FromString(templateStr)
	if err != nil {
		panic(err)
	}
	out, err := tmpl.Execute(ctx)
	if err != nil {
		panic(err)
	}

	fmt.Println(out)
}
