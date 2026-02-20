package main

import (
	"encoding/json"
	"log"
	"os"
	"time"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

func main() {
	var messages []int
	seen := map[int]bool{}

	n := maelstrom.NewNode()

	n.Handle("broadcast", func(msg maelstrom.Message) error {
		var request map[string]any
		if err := json.Unmarshal(msg.Body, &request); err != nil {
			return err
		}

		response := map[string]any{}
		response["type"] = "broadcast_ok"

		message := int(request["message"].(float64))
		if seen[message] {
			return nil
		}
		messages = append(messages, message)
		seen[message] = true

		for _, peer := range n.NodeIDs() {
			if peer != n.ID() {
				n.Send(peer, request)
			}
		}

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

		_ = topologyRequest.Topology

		response := map[string]any{}
		response["type"] = "topology_ok"

		return n.Reply(msg, response)
	})

	ticker := time.NewTicker(500 * time.Millisecond)
	go func() {
		for range ticker.C {
			for _, peer := range n.NodeIDs() {
				if peer != n.ID() {
					for _, number := range messages {
						n.RPC(peer, broadcastRequest(number), nil)
					}

				}
			}
		}
	}()

	if err := n.Run(); err != nil {
		log.Printf("ERROR: %s", err)
		os.Exit(1)
	}
}

func broadcastRequest(number int) map[string]any {
	request := map[string]any{}

	request["type"] = "broadcast"
	request["message"] = number

	return request
}
