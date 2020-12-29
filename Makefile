start-minikube-with-calico:
	minikube start --cni=calico

deploy:
	kustomize build scripts/example | kubectl apply -f -