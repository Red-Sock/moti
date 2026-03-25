build-bin:
	go build -o $(GOBIN)/moti main.go

lint:
	golangci-lint run