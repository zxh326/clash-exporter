NAME=clash-exporter
VERSION=$(shell git describe --tags --always)

releases: darwin-amd64 darwin-arm64 linux-amd64 linux-arm64 linux-armv6 linux-armv7

darwin-amd64:
	GOOS=darwin GOARCH=amd64 go build -o $(NAME)-${VERSION}-$@
	tar czf $(NAME)-${VERSION}-$@.tar.gz $(NAME)-${VERSION}-$@
darwin-arm64:
	GOOS=darwin GOARCH=arm64 go build -o $(NAME)-${VERSION}-$@
	tar czf $(NAME)-${VERSION}-$@.tar.gz $(NAME)-${VERSION}-$@
linux-amd64:
	GOOS=linux GOARCH=amd64 go build -o $(NAME)-${VERSION}-$@
	tar czf $(NAME)-${VERSION}-$@.tar.gz $(NAME)-${VERSION}-$@
linux-arm64:
	GOOS=linux GOARCH=arm64 go build -o $(NAME)-${VERSION}-$@
	tar czf $(NAME)-${VERSION}-$@.tar.gz $(NAME)-${VERSION}-$@
linux-armv6:
	GOOS=linux GOARCH=arm GOARM=6 go build -o $(NAME)-${VERSION}-$@
	tar czf $(NAME)-${VERSION}-$@.tar.gz $(NAME)-${VERSION}-$@
linux-armv7:
	GOOS=linux GOARCH=arm GOARM=7 go build -o $(NAME)-${VERSION}-$@
	tar czf $(NAME)-${VERSION}-$@.tar.gz $(NAME)-${VERSION}-$@

clean:
	rm -rf $(NAME)-*