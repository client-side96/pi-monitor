build:
	env GOOS=linux GOARCH=arm GOARM=5 go build  -o bin/pi-monitor ./cmd/pi-monitor/main.go

test:
	go test -v ./...

lint:
	golangci-lint run

image:
	make build
	docker build -t pi-monitor:local .

start:
	docker run -d -p 4000:4000 --rm --platform linux/arm/v5 pi-monitor:local

