package main

import (
	"encoding/json"
	"log"
	"os"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

func main() {
	n := maelstrom.NewNode()

	// send: append a message to a log key, return the assigned offset.
	//
	// Request:  {"type": "send", "key": "k1", "msg": 123}
	// Response: {"type": "send_ok", "offset": <int>}
	//
	// Offsets are monotonically increasing per key (may be sparse).
	n.Handle("send", func(msg maelstrom.Message) error {
		var req struct {
			Key string `json:"key"`
			Msg int    `json:"msg"`
		}
		if err := json.Unmarshal(msg.Body, &req); err != nil {
			return err
		}

		// TODO(human): append req.Msg to the log for req.Key,
		// assign an offset, and reply with {"type": "send_ok", "offset": offset}.
		_ = req

		return n.Reply(msg, map[string]any{"type": "send_ok", "offset": 0})
	})

	// poll: return messages from each requested log starting at the given offset.
	//
	// Request:  {"type": "poll", "offsets": {"k1": 1000, "k2": 2000}}
	// Response: {"type": "poll_ok", "msgs": {"k1": [[1000, 9], [1001, 5]], "k2": [[2000, 7]]}}
	//
	// Each entry is [offset, value]. Return messages starting from the requested offset.
	n.Handle("poll", func(msg maelstrom.Message) error {
		var req struct {
			Offsets map[string]int `json:"offsets"`
		}
		if err := json.Unmarshal(msg.Body, &req); err != nil {
			return err
		}

		// TODO(human): for each key in req.Offsets, return messages from that offset onward.
		// msgs is map[string][][2]int — each entry is [offset, value].
		_ = req

		return n.Reply(msg, map[string]any{"type": "poll_ok", "msgs": map[string]any{}})
	})

	// commit_offsets: record that messages up to the given offset have been processed.
	//
	// Request:  {"type": "commit_offsets", "offsets": {"k1": 1000, "k2": 2000}}
	// Response: {"type": "commit_offsets_ok"}
	n.Handle("commit_offsets", func(msg maelstrom.Message) error {
		var req struct {
			Offsets map[string]int `json:"offsets"`
		}
		if err := json.Unmarshal(msg.Body, &req); err != nil {
			return err
		}

		// TODO(human): store the committed offsets for each key.
		_ = req

		return n.Reply(msg, map[string]any{"type": "commit_offsets_ok"})
	})

	// list_committed_offsets: return the last committed offset for each requested key.
	//
	// Request:  {"type": "list_committed_offsets", "keys": ["k1", "k2"]}
	// Response: {"type": "list_committed_offsets_ok", "offsets": {"k1": 1000, "k2": 2000}}
	n.Handle("list_committed_offsets", func(msg maelstrom.Message) error {
		var req struct {
			Keys []string `json:"keys"`
		}
		if err := json.Unmarshal(msg.Body, &req); err != nil {
			return err
		}

		// TODO(human): look up the committed offset for each key in req.Keys.
		_ = req

		return n.Reply(msg, map[string]any{"type": "list_committed_offsets_ok", "offsets": map[string]int{}})
	})

	if err := n.Run(); err != nil {
		log.Printf("ERROR: %s", err)
		os.Exit(1)
	}
}
