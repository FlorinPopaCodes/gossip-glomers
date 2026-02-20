package main

import (
	"encoding/json"
	"log"
	"os"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

func main() {
	var messages []int

	n := maelstrom.NewNode()

	n.Handle("broadcast", func(msg maelstrom.Message) error {
		var request map[string]any
		response := map[string]any{}
		if err := json.Unmarshal(msg.Body, &request); err != nil {
			return err
		}

		messages = append(messages, int(request["message"].(float64)))

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
		response := map[string]any{}

		response["type"] = "topology_ok"

		return n.Reply(msg, response)
	})

	if err := n.Run(); err != nil {
		log.Printf("ERROR: %s", err)
		os.Exit(1)
	}
}
