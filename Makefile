buid:
	go build rest2mail.go
docker:
	go test && go build -o main .
	docker build -t frehub/rest2mail .