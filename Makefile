    .PHONY: all fmt lint build test unit integration docker compose-up compose-down kind-up deploy e2e smoke clean zip

    APP?=ordersvc
    GOFLAGS=-trimpath
    LDFLAGS=-s -w
    BIN_DIR=bin

    all: fmt build test

    fmt:
	go fmt ./...

    build:
	mkdir -p $(BIN_DIR)
	GOFLAGS=$(GOFLAGS) go build -ldflags "$(LDFLAGS)" -o $(BIN_DIR)/ordersvc ./cmd/ordersvc
	GOFLAGS=$(GOFLAGS) go build -ldflags "$(LDFLAGS)" -o $(BIN_DIR)/paymentsvc ./cmd/paymentsvc

    test: unit integration

    unit:
	go test ./... -run Test -count=1

    integration:
	go test ./test/integration -tags=integration -count=1 || echo "Integration tests require local env; see README"

    docker:
	docker build -t wli/ordersvc -f deploy/docker/ordersvc.Dockerfile .
	docker build -t wli/paymentsvc -f deploy/docker/paymentsvc.Dockerfile .

    compose-up:
	docker compose up --build -d

    compose-down:
	docker compose down -v

    kind-up:
	scripts/kind-up.sh

    deploy:
	kubectl apply -k deploy/kustomize

    e2e:
	scripts/e2e.sh

    smoke:
	scripts/smoke.sh

    zip:
	python3 scripts/make_zip.py

    clean:
	rm -rf $(BIN_DIR)
	git clean -fdX || true
