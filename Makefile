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
	kubectl create namespace gpsd || true

setup:
	kubectl create configmap gpsd-nginx-config --from-file=nginx.conf -n gpsd || kubectl replace configmap gpsd-nginx-config --from-file=nginx.conf -n gpsd
	kubectl apply -f deployments/api-gateway-deployment.yaml
	kubectl apply -f deployments/nginx-deployment.yaml
	kubectl apply -f services/api-gateway-service.yaml
	kubectl apply -f services/nginx-service.yaml

run:
	sleep 5
	kubectl get pods -n gpsd
	kubectl get services -n gpsd
	minikube service gpsd-nginx-service -n gpsd

all: build build-image setup run

clean:
	kubectl delete all --all -n gpsd || true
	kubectl delete configmap gpsd-nginx-config -n gpsd

	kubectl delete namespace gpsd || true
	sleep 2