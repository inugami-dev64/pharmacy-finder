vpath % ./frontend/node_modules ./frontend/build 

GOSRC := $(shell find . -type f -name "*.go" -o -name "*.sql" 2>/dev/null)

./frontend/node_modules: ./frontend/package.json
	npm install --prefix ./frontend

./frontend/build: ./frontend/node_modules
	npm --prefix frontend run build

./pharmafinder-dbg: go.mod go.sum ${GOSRC}
	CGO_ENABLED=0 go build -o pharmafinder-dbg ./cmd/pharmafinder/main.go

./pharmafinder: go.mod go.sum ${GOSRC}
	go build -o pharmafinder -ldflags="-s -w" ./cmd/pharmafinder/main.go

.PHONY: debug
debug: ./frontend/build ./pharmafinder-dbg

.PHONY: release
release: ./frontend/build ./pharmafinder
	
.PHONY: clean
clean:
	rm -rf ./frontend/node_modules
	rm -rf ./frontend/build
	rm pharmafinder*