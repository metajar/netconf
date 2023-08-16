package rpcs

import (
	"fmt"
	"github.com/metajar/netconf/netconf"
)

func DefaultLogRpcReplyCallback(eventId string) netconf.Callback {
	return func(event netconf.Event) {
		reply := event.RPCReply()
		if reply == nil {
			println("Failed to execute RPC")
		}
		if event.EventID() == eventId {
			println("Successfully executed RPC")
			println(reply.RawReply)
		}
	}
}

func DefaultGatherer(eventId string, payload *string) netconf.Callback {
	return func(event netconf.Event) {
		fmt.Println("Doing the thing")
		reply := event.RPCReply()
		fmt.Println(reply.Data)
		if reply == nil {
			println("Failed to execute RPC")
		}
		if event.EventID() == eventId {
			println("Successfully executed RPC")
			payload = &reply.RawReply
		}
	}
}
