.PHONY: build test lint typecheck security e2e check

build:
	cd backend && go build -o bin/kubemv ./cmd/kubemv
	cd desktop && npm run build
	cd desktop/src-tauri && cargo check

test:
	cd backend && go test ./... -count=1
	cd desktop && npm test

lint:
	cd backend && test -z "$$(gofmt -l .)"
	cd backend && go vet ./...
	cd desktop && npm run lint

typecheck:
	cd desktop && npm run typecheck

security:
	./scripts/security-check.sh

e2e:
	cd desktop && npm run e2e

check: lint typecheck test security build
