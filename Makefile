build:
	go build -o ./rain -v ./main.go

run:
	./rain -loglevel 0 

.DEFAULT_GOAL := build
