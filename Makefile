build:
	GOOS=js GOARCH=wasm go build -o docs/main.wasm .
	cp `go env GOROOT`/misc/wasm/wasm_exec.js docs/

run:
	python -m http.server -d docs
