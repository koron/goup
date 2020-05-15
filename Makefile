.PHONY: build
build:
	go build

.PHONY: test
test:
	go test

.PHONY: tags
tags:
	gotags -f tags -R .
