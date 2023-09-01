build:
	GOARCH=arm64 GOOS=linux go build -o ./out/collect_launches/bootstrap -ldflags "-s -w" ./functions/collectNextRocketLaunches
test:
	go test ./functions/...