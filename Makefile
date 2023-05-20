BIN=yacht_dice

build:
	go build -o $(BIN) .

clean:
	rm $(BIN)
