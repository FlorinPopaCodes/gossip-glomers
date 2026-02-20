CHALLENGES = maelstrom-echo maelstrom-unique-ids

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
