CHALLENGES = maelstrom-echo maelstrom-unique-ids maelstrom-broadcast maelstrom-counter

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

test-broadcast-3c: maelstrom-broadcast
	maelstrom test -w broadcast \
		--bin ./bin/maelstrom-broadcast \
		--node-count 5 \
		--time-limit 20 \
		--rate 10 \
		--nemesis partition

test-broadcast-3d: maelstrom-broadcast
	maelstrom test -w broadcast \
		--bin ./bin/maelstrom-broadcast \
		--node-count 25 \
		--time-limit 20 \
		--rate 100 \
		--latency 100
	./check-results.sh 30 400 600

test-broadcast: maelstrom-broadcast
	maelstrom test -w broadcast \
		--bin ./bin/maelstrom-broadcast \
		--node-count 25 \
		--time-limit 20 \
		--rate 100 \
		--latency 100
	./check-results.sh 20 1000 2000

test-counter: maelstrom-counter
	maelstrom test -w g-counter \
		--bin ./bin/maelstrom-counter \
		--node-count 3 \
		--rate 100 \
		--time-limit 20 \
		--nemesis partition
