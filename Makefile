build:
	go build -o bin/main .

run: build
	./bin/main

test-actors:
	go test -timeout 300s -run ^TestActors$