NAMESPACE = gpsd
DEPLOYMENT = gpsd-api-gateway
IMAGE_NAME = $(NAMESPACE)/$(DEPLOYMENT)
TAG ?= latest  # If no tag is provided, default to 'latest'

develop: helm-uninstall build-image push-image helm

all: build build-image push-image setup run

build-image:
	docker build -f Dockerfile -t $(IMAGE_NAME):$(TAG) .

push-image:
	docker push $(IMAGE_NAME):$(TAG)

run-image:
	docker run -p 3000:3000 $(DEPLOYMENT)

clean-image:
	docker rmi $(docker images --filter "dangling=true" -q) -f

certs:
	chmod +x private/pki.sh
	bash private/pki.sh

helm:
	helm upgrade --install demo ./stage1 --set image.tag=$(TAG) --namespace $(NAMESPACE)

helm-uninstall:
	helm uninstall demo -n $(NAMESPACE) 

clean:
	kubectl delete all --all -n $(NAMESPACE)  || true
	kubectl delete namespace $(NAMESPACE)  || true
	sleep 2