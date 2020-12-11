export GO111MODULE = on
export GOPROXY = https://proxy.golang.org

.PHONE: build
build:
	cp "${GOROOT}/misc/wasm/wasm_exec.js" ./lib
	GOOS=js GOARCH=wasm go build -trimpath --ldflags "-s -w" -v -o dist/go.wasm main.go

publish:
	curl -X PUT \
		"https://api.cloudflare.com/client/v4/accounts/${CF_ACCOUNT_ID}/workers/scripts/${CF_WORKER_NAME}" \
		-H "Authorization: Bearer ${CF_API_TOKEN}" \
		-F "metadata=@dist/metadata.json;type=application/json" \
		-F "script=@dist/worker.js;type=application/javascript" \
		-F "wasm=@dist/go.wasm;type=application/wasm"

run:
	GOOS=js GOARCH=wasm go run main.go

fmt:
	@echo "-> Running go fmt"
	@go fmt ./...

clean:
	rm -rf .cache dist
