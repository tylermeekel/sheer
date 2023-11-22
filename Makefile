build:
	go build -o ./bin

send: build
	./bin/p2pfileshare s ./test.txt

receive: build
	./bin/p2pfileshare r
