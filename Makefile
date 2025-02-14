all: build build-image push-image setup run

build-image:
	docker build -f Dockerfile -t gpsd/gpsd-api-gateway:v1 .

push-image:
	docker push gpsd/gpsd-api-gateway:v1

run-image:
	docker run -p 3000:3000 gpsd-api-gateway

clean-image:
	docker rmi $(docker images --filter "dangling=true" -q) -f

certs:
	chmod +x private/pki.sh
	bash private/pki.sh

build:
	kubectl create namespace gpsd || true

setup:
	kubectl create secret generic gpsd-api-gateway-secrets --from-literal=JWT_SECRET=gpsdjwtsecretkey --from-literal=REFRESH_SECRET=gpsdrefreshsecretkey -n gpsd
	kubectl create secret tls gpsd-api-gateway-certificates --cert=private/certs/api.gpsd.com.crt --key=private/certs/api.gpsd.com.key -n gpsd
	kubectl apply -f deployments/api-gateway-deployment.yaml
	kubectl apply -f services/api-gateway-service.yaml

run:
	sleep 5
	kubectl get pods -n gpsd
	kubectl get services -n gpsd

clean:
	kubectl delete all --all -n gpsd || true
	kubectl delete secret gpsd-api-gateway-secrets -n gpsd
	kubectl delete secret gpsd-api-gateway-certificates -n gpsd

	kubectl delete namespace gpsd || true
	sleep 2