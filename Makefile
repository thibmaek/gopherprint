include Makefile.vars

LDFLAGS="-X 'main.apiKey=$(API_KEY)' -X 'main.serverURL=$(SERVER_URL)'"

install_deps:
	go get github.com/caseymrm/menuet
	go get github.com/mcuadros/go-octoprint
	go get github.com/haklop/gnotifier

build:
	rm -rf Gopherprint.app
	go build -ldflags=$(LDFLAGS) -o appTemplate/Contents/MacOS/gopherprint gopherprint.go
	cp -r appTemplate Gopherprint.app

install: install_deps build
	cp -r Gopherprint.app /Applications/
