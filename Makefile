start-minikube-with-calico:
	minikube start --cni=calico

deploy:
	kustomize build scripts/example | kubectl apply -f -

run:
	go run cmd/main.go

test:
	go test ./... -count=5 -race