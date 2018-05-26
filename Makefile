test:
	go test . ./actions

show:
	kubectl get all,secrets,configmaps,certificates,ingress
	gcloud compute addresses list

logs:
	kubectl logs `kubectl get pods -l 'app=faas,tier=web' -o jsonpath="{.items[0].metadata.name}"` -c faas-app -f

build:
	buffalo build -o bin/faas

upload:
	docker build -t gcr.io/faas-203916/faas:latest .
	docker push gcr.io/faas-203916/faas:latest

restart:
	kubectl delete pod/`kubectl get pods -l 'app=faas,tier=web' -o jsonpath="{.items[0].metadata.name}"`

update: build upload restart
