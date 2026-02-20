CHALLENGES = maelstrom-echo maelstrom-unique-ids maelstrom-broadcast

.PHONY: all clean $(CHALLENGES)

all: $(CHALLENGES)

$(CHALLENGES):
	go build -o ./bin/$@ ./cmd/$@

clean:
	rm -rf ./bin

# --- Test targets ---

test-echo: maelstrom-echo
	maelstrom test -w echo \
		--bin ./bin/maelstrom-echo \
		--node-count 1 \
		--time-limit 10

test-unique-ids: maelstrom-unique-ids
	maelstrom test -w unique-ids \
		--bin ./bin/maelstrom-unique-ids \
		--time-limit 30 \
		--rate 1000 \
		--node-count 3 \
		--availability total \
		--nemesis partition

test-broadcast-3a: maelstrom-broadcast
	maelstrom test -w broadcast \
		--bin ./bin/maelstrom-broadcast \
		--node-count 1 \
		--time-limit 20 \
		--rate 10

test-broadcast-3b: maelstrom-broadcast
	maelstrom test -w broadcast \
		--bin ./bin/maelstrom-broadcast \
		--node-count 5 \
		--time-limit 20 \
		--rate 10

test-broadcast: maelstrom-broadcast
	maelstrom test -w broadcast \
		--bin ./bin/maelstrom-broadcast \
		--node-count 5 \
		--time-limit 20 \
		--rate 10 \
		--nemesis partition
