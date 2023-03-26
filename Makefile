.DEFAULT_TARGET = build

.SLONY = build
build:
	go build -o ./bin/actors .

.SLONY = run
run: build
	./bin/actors

.SLONY = test
test: 
	go test ./...