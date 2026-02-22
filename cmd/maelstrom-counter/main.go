package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"os"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

func main() {
	n := maelstrom.NewNode()
	kv := maelstrom.NewSeqKV(n)
	ctx := context.Background()

	n.Handle("add", func(msg maelstrom.Message) error {
		var addRequest struct {
			Delta int `json:"delta"`
		}
		if err := json.Unmarshal(msg.Body, &addRequest); err != nil {
			return err
		}

		currentValue, err := readIntOr(kv, ctx, n.ID(), 0)
		if err != nil {
			return err
		}

		newValue := currentValue + addRequest.Delta

		if err := kv.CompareAndSwap(ctx, n.ID(), currentValue, newValue, true); err != nil {
			return err
		}

		return n.Reply(msg, map[string]any{"type": "add_ok"})
	})

	n.Handle("read", func(msg maelstrom.Message) error {
		var result int

		for _, nodeID := range n.NodeIDs() {
			val, err := readIntOr(kv, ctx, nodeID, 0)
			if err != nil {
				return err
			}

			result += val
		}

		return n.Reply(msg, map[string]any{
			"type":  "read_ok",
			"value": result,
		})
	})

	if err := n.Run(); err != nil {
		log.Printf("ERROR: %s", err)
		os.Exit(1)
	}

}

func readIntOr(kv *maelstrom.KV, ctx context.Context, key string, fallback int) (int, error) {
	val, err := kv.ReadInt(ctx, key)
	if err != nil {
		var rpcErr *maelstrom.RPCError
		if errors.As(err, &rpcErr) && rpcErr.Code == maelstrom.KeyDoesNotExist {
			return fallback, nil
		}
		return 0, err
	}
	return val, nil
}
