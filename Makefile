include Makefile.vars

LDFLAGS="-X 'main.apiKey=$(API_KEY)' -X 'main.serverURL=$(SERVER_URL)'"

build:
	rm -rf Gopherprint.app
	go build -ldflags=$(LDFLAGS) -o appTemplate/Contents/MacOS/gopherprint gopherprint.go
	cp -r appTemplate Gopherprint.app

install: build
	cp -r Gopherprint.app /Applications/
