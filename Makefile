wasm:
	GOARCH=wasm GOOS=js go build -ldflags="-s -w" -trimpath -o _web/main.wasm ./cmd/game

itchio-wasm: wasm
	cd _web && \
		mkdir -p ../bin && \
		rm -f ../bin/pixelspace_rangers.zip && \
		zip ../bin/pixelspace_rangers.zip -r main.wasm index.html wasm_exec.js
