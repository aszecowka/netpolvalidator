export GOBIN=$(PWD)/bin/go

.DEFAULT_GOAL = all

start-minikube-with-calico:
	minikube start --cni=calico

deploy:
	kustomize build scripts/example | kubectl apply -f -

run:
	go run cmd/main.go

test:
	go test ./... -count=5 -race

build:
	go build -o ./bin/netpol cmd/main.go

check-dependencies:
	go mod verify
	go mod tidy -v

fmt:
	go fmt ./...

static-checks:
	go install github.com/kisielk/errcheck
	$(GOBIN)/errcheck -blank ./...
	go vet ./...


all: check-dependencies build test fmt static-checks