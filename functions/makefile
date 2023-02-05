build-all:
	GOOS=linux GOARCH=amd64
	cd collectPeopleInSpace && go build -o ../../out/collect_people_linux -ldflags "-s -w" && cd ..