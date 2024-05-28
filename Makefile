build:
	env GOOS=linux GOARCH=arm GOARM=5 go build  -o bin/pi-monitor ./cmd/pi-monitor/main.go

test:
	go test -v ./...

lint:
	golangci-lint run

image:
	make build
	docker buildx build --platform linux/arm/v5 . -t pi-monitor:local

start:
	docker run -d -p 8000:8000 --rm --platform linux/arm/v5 pi-monitor:local

