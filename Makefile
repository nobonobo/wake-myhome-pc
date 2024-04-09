OUTPUT=$(shell find docs | grep "docs/[0-9a-f]\+$$")

build:
	GOOS=js GOARCH=wasm go build -o $(OUTPUT)/main.wasm .
	cp `go env GOROOT`/misc/wasm/wasm_exec.js $(OUTPUT)/

run:
	python -m http.server -d docs
