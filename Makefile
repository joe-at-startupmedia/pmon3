build:
	go mod tidy
	go build -o bin/pmon2 cmd/pmon2/pmon2.go
	go build -o bin/pmond cmd/pmond/pmond.go
install:
	rm -rf /usr/local/pmon2/bin/
	cp -R bin/ /usr/local/pmon2/
