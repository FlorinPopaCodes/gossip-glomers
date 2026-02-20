package main

import (
	"encoding/json"
	"log"
	"os"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

func main() {
	var messages []int
	topology := map[string][]string{}

	n := maelstrom.NewNode()

	n.Handle("broadcast", func(msg maelstrom.Message) error {
		var request map[string]any
		if err := json.Unmarshal(msg.Body, &request); err != nil {
			return err
		}

		message := int(request["message"].(float64))
		messages = append(messages, message)

		if topology[msg.Src] == nil {
			for _, peer := range n.NodeIDs() {
				if peer != n.ID() {
					n.Send(peer, request)
				}
			}
		}

		response := map[string]any{}
		response["type"] = "broadcast_ok"

		return n.Reply(msg, response)
	})

	n.Handle("read", func(msg maelstrom.Message) error {
		response := map[string]any{}

		response["type"] = "read_ok"
		response["messages"] = messages

		return n.Reply(msg, response)
	})

	n.Handle("topology", func(msg maelstrom.Message) error {
		var topologyRequest struct {
			Topology map[string][]string `json:"topology"`
		}
		if err := json.Unmarshal(msg.Body, &topologyRequest); err != nil {
			return err
		}

		topology = topologyRequest.Topology

		response := map[string]any{}
		response["type"] = "topology_ok"

		return n.Reply(msg, response)
	})

	if err := n.Run(); err != nil {
		log.Printf("ERROR: %s", err)
		os.Exit(1)
	}
}
