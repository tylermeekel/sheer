build:
	go build -o ./bin/sheer

send: build
	./bin/sheer s ./test.txt

receive: build
	./bin/sheer r
