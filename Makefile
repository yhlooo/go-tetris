build-wasm:
	GOARCH=wasm GOOS=js go build -o web/app.wasm ./cmd/tetris-wasm
	go build -o web/server ./cmd/tetris-wasm

run-web: build-wasm
	./web/server
