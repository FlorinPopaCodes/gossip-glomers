# Gossip Glomers

Solutions to the [Fly.io distributed systems challenges](https://fly.io/dist-sys/) using [Maelstrom](https://github.com/jepsen-io/maelstrom).

## Challenges

| # | Challenge | Binary |
|---|-----------|--------|
| 1 | Echo | `maelstrom-echo` |
| 2 | Unique ID Generation | `maelstrom-unique-ids` |
| 3 | Broadcast (3/5) | `maelstrom-broadcast` |
| 4 | Grow-Only Counter | `maelstrom-counter` |
| 5 | Kafka-Style Log | `maelstrom-kafka` |
| 6 | Totally-Available Transactions | `maelstrom-txn` |

## Build & Test

```bash
make all              # build all challenges
make test-echo        # run Maelstrom echo test
make test-unique-ids  # run Maelstrom unique-ids test
make test-broadcast-3a # single-node broadcast
make test-broadcast-3b # multi-node broadcast
make test-broadcast-3c # fault tolerant broadcast
make test-counter     # grow-only counter with partitions
make test-kafka       # single-node kafka-style log
make test-txn         # single-node transactions
```

Requires [Maelstrom](https://github.com/jepsen-io/maelstrom) in your PATH.

## License

[MIT](LICENSE)
