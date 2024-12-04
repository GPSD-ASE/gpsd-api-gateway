build-image:
	docker build -f docker/Dockerfile.api -t paranjik/gpsd-api-gateway:latest .
	docker build -f docker/Dockerfile.nginx -t paranjik/gpsd-nginx:latest .

push-image:
	docker push paranjik/gpsd-api-gateway:latest
	docker push paranjik/gpsd-nginx:latest

run-image:
	docker run -p 3000:3000 gpsd-api-gateway
	docker run -p 80:80 gpsd-nginx

clean-image:
	docker rmi $(docker images --filter "dangling=true" -q) -f

build:
	make build-image
	kubectl create namespace gpsd
	kubectl apply -f deployments/api-gateway-deployment.yaml
	kubectl apply -f deployments/nginx-deployment.yaml
	kubectl apply -f services/api-gateway-service.yaml
	kubectl apply -f services/nginx-service.yaml

run:
	sleep 2
	minikube service gpsd-nginx-service -n gpsd
	chmod +x commands.sh
	bash commands.sh

all: build run

clean:
	kubectl delete all --all -n gpsd
	kubectl delete namespace gpsd
	make clean-image