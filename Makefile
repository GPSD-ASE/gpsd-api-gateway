TAG ?= 0.1.2  # If no tag is provided, default to 'latest'
NAMESPACE = gpsd
DEPLOYMENT = gpsd-api-gateway
SERVICE_NAME = $(DEPLOYMENT)
IMAGE_NAME = $(NAMESPACE)/$(DEPLOYMENT)
CHART_DIRECTORY = helm
LOCAL_CHART_NAME = $(shell ls /tmp/$(DEPLOYMENT)-*.tgz)
LOCAL_INDEX_FILE = /tmp/index.yaml
REMOTE_CHART_REPOSITORY = gpsd-ase.github.io

# Use `make develop` for local testing
develop: helm-uninstall build-image push-image helm

docker: build-image push-image

build-image:
	docker build -f Dockerfile -t $(IMAGE_NAME):$(TAG) --platform linux/amd64 .

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
	@echo "Upgrading/Installing $(DEPLOYMENT) Helm chart..."
	helm upgrade --install $(DEPLOYMENT) ./helm --set image.tag=$(TAG) --namespace $(NAMESPACE)

helm-uninstall:
	@echo "Uninstalling $(DEPLOYMENT) from Kubernetes..."
	helm uninstall demo -n $(NAMESPACE) 

clean:
	@echo "Cleaning up all resources in the $(NAMESPACE) namespace..."
	kubectl delete all --all -n $(NAMESPACE)  || true
	kubectl delete namespace $(NAMESPACE)  || true
	sleep 2


deploy-gh-pages: gh-pages-publish helm-repo-update

gh-pages-publish:
	@echo "Publishing Helm chart for $(SERVICE_NAME) to GitHub Pages..."
	rm -rf $(LOCAL_CHART_NAME) $(LOCAL_INDEX_FILE) /tmp/$(NAMESPACE)-*.tgz
	helm package ./$(CHART_DIRECTORY) -d /tmp
	helm repo index /tmp --url https://$(REMOTE_CHART_REPOSITORY)/$(SERVICE_NAME)/ --merge /tmp/index.yaml
	git checkout gh-pages
	cp  $(LOCAL_CHART_NAME) $(LOCAL_INDEX_FILE) .
	git add .
	git commit -m "fix: commit to update GitHub Pages"
	git push origin gh-pages -f
	watch curl -k https://$(REMOTE_CHART_REPOSITORY)/$(SERVICE_NAME)/index.yaml

helm-repo-update:
	@echo "Adding and updating Helm repo for $(SERVICE_NAME)..."
	helm repo add $(SERVICE_NAME) https://$(REMOTE_CHART_REPOSITORY)/$(SERVICE_NAME)/
	helm repo update
	helm repo list