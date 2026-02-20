# Gossip Glomers

Solutions to the [Fly.io distributed systems challenges](https://fly.io/dist-sys/) using [Maelstrom](https://github.com/jepsen-io/maelstrom).

## Challenges

| # | Challenge | Binary |
|---|-----------|--------|
| 1 | Echo | `maelstrom-echo` |
| 2 | Unique ID Generation | `maelstrom-unique-ids` |
| 3 | Broadcast | _todo_ |
| 4 | Grow-Only Counter | _todo_ |
| 5 | Kafka-Style Log | _todo_ |
| 6 | Totally-Available Transactions | _todo_ |

## Build & Test

```bash
make all              # build all challenges
make test-echo        # run Maelstrom echo test
make test-unique-ids  # run Maelstrom unique-ids test
```

Requires [Maelstrom](https://github.com/jepsen-io/maelstrom) in your PATH.

## License

[MIT](LICENSE)
