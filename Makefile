SERVER_URL=$(url)
API_KEY=$(key)
LDFLAGS="-X 'main.apiKey=$(API_KEY)' -X 'main.serverURL=$(SERVER_URL)'"

build:
	rm -rf Gopherprint.app
	go build -ldflags=$(LDFLAGS) -o appTemplate/Contents/MacOS/gopherprint gopherprint.go
	cp -r appTemplate Gopherprint.app
