build:
	cd collectPeopleInSpace && GOARCH=amd64 GOOS=linux go build -o ../../out/collect_people_linux -ldflags "-s -w" && cd ..
	cd readPeopleInSpace && GOARCH=amd64 GOOS=linux go build -o ../../out/read_people_linux -ldflags "-s -w" && cd ..
	cd collectNextRocketLaunches && GOARCH=amd64 GOOS=linux go build -o ../../out/collect_rockets_linux -ldflags "-s -w" && cd ..