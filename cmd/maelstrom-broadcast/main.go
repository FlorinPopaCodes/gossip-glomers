package main

import (
	"encoding/json"
	"log"
	"os"
	"sync"
	"time"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

func main() {
	var mu sync.RWMutex
	var messages []int
	seen := map[int]struct{}{}
	var siblings []string

	n := maelstrom.NewNode()

	n.Handle("broadcast", func(msg maelstrom.Message) error {
		var request map[string]any
		if err := json.Unmarshal(msg.Body, &request); err != nil {
			return err
		}

		response := map[string]any{}
		response["type"] = "broadcast_ok"

		message := int(request["message"].(float64))
		mu.Lock()
		if _, exists := seen[message]; exists {
			mu.Unlock()
			return nil
		}
		messages = append(messages, message)
		seen[message] = struct{}{}
		mu.Unlock()

		// TODO: remove this for part e
		for _, peer := range siblings {
			if peer != msg.Src {
				// TODO: start tracking acks
				n.Send(peer, request)
			}
		}

		return n.Reply(msg, response)
	})

	n.Handle("read", func(msg maelstrom.Message) error {
		response := map[string]any{}

		response["type"] = "read_ok"

		mu.RLock()
		snapshot := make([]int, len(messages))
		copy(snapshot, messages)
		mu.RUnlock()

		response["messages"] = snapshot

		return n.Reply(msg, response)
	})

	n.Handle("topology", func(msg maelstrom.Message) error {
		var topologyRequest struct {
			Topology map[string][]string `json:"topology"`
		}
		if err := json.Unmarshal(msg.Body, &topologyRequest); err != nil {
			return err
		}

		// TODO: Remove loops from siblings
		siblings = topologyRequest.Topology[n.ID()]

		response := map[string]any{}
		response["type"] = "topology_ok"

		return n.Reply(msg, response)
	})

	ticker := time.NewTicker(1 * time.Second)
	go func() {
		for range ticker.C {
			// TODO: only re-send unacks
			mu.RLock()
			snapshot := make([]int, len(messages))
			copy(snapshot, messages)
			mu.RUnlock()

			for _, peer := range siblings {
				for _, number := range snapshot {
					n.RPC(peer, broadcastRequest(number), nil)
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
