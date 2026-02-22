package main

import (
	"encoding/json"
	"log"
	"os"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

func main() {
	n := maelstrom.NewNode()

	// txn: execute a list of read/write operations atomically.
	//
	// Request:  {"type": "txn", "txn": [["r", 1, null], ["w", 1, 6], ["w", 2, 9]]}
	// Response: {"type": "txn_ok", "txn": [["r", 1, 3], ["w", 1, 6], ["w", 2, 9]]}
	//
	// Operations are [op, key, value] triples:
	//   ["r", key, null]  — read: replace null with the current value (or null if missing)
	//   ["w", key, value] — write: set key to value
	n.Handle("txn", func(msg maelstrom.Message) error {
		var req struct {
			Txn [][]any `json:"txn"`
		}
		if err := json.Unmarshal(msg.Body, &req); err != nil {
			return err
		}

		// TODO(human): process each operation in req.Txn.
		// For reads: look up the key and fill in the value (index 2).
		// For writes: store the value for the key.
		// Return the updated txn array in the response.
		_ = req

		return n.Reply(msg, map[string]any{"type": "txn_ok", "txn": req.Txn})
	})

	if err := n.Run(); err != nil {
		log.Printf("ERROR: %s", err)
		os.Exit(1)
	}
}
