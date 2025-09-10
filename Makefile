vpath % ./frontend/node_modules ./frontend/build

GOPATH := $(shell go env GOPATH)
GOSRC := $(shell find . -type f -name "*.go" ! -regex '.*_mock.go' -o -name "*.sql" 2>/dev/null)

MOCKGEN_DST := mock/pharmacy_repository_mock.go \
			   mock/http_mock.go \
			   mock/db_mock.go

.PHONY: all
all: ./pharmafinder-dbg ./pharmafinder

./frontend/node_modules: ./frontend/package.json
	npm install --prefix ./frontend

./frontend/build: ./frontend/node_modules
	npm --prefix frontend run build

./pharmafinder-dbg: go.mod go.sum ${GOSRC}
	CGO_ENABLED=0 go build -gcflags="-N -l" -o pharmafinder-dbg ./cmd/pharmafinder/main.go

./pharmafinder: go.mod go.sum ${GOSRC}
	go build -o pharmafinder -ldflags="-s -w" ./cmd/pharmafinder/main.go

.PHONY: debug
debug: ./frontend/build ./pharmafinder-dbg

.PHONY: release
release: ./frontend/build ./pharmafinder

### Mockgen targets ###
# Repositories
mock/pharmacy_repository_mock.go: db/pharmacy_repository.go
	${GOPATH}/bin/mockgen -source=db/pharmacy_repository.go -destination=mock/pharmacy_repository_mock.go -package=mock

mock/http_mock.go: utils/http.go
	${GOPATH}/bin/mockgen -source=utils/http.go -destination=mock/http_mock.go -package=mock

mock/db_mock.go: db/db.go
	${GOPATH}/bin/mockgen -source=db/db.go -destination=mock/db_mock.go -package=mock

.PHONY: mockgen
mockgen: ${MOCKGEN_DST}

.PHONY: clean
clean:
	rm -rf mock
	rm -rf ./frontend/node_modules
	rm -rf ./frontend/build
	rm pharmafinder*