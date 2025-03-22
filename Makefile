TAG ?= 0.1.3
NAMESPACE = gpsd
DEPLOYMENT = gpsd-api-gateway
SERVICE_NAME = $(DEPLOYMENT)
IMAGE_NAME = $(NAMESPACE)/$(DEPLOYMENT)
CHART_DIRECTORY = helm
REMOTE_CHART_REPOSITORY = gpsd-ase.github.io
VERSION := $(shell grep "version:" helm/Chart.yaml | head -1 | sed 's/version: //')

# Docker commands
.PHONY: docker build-image push-image run-image clean-image
docker: build-image push-image

build-image:
	@echo "Building Docker image $(IMAGE_NAME):$(TAG)..."
	docker build -f Dockerfile -t $(IMAGE_NAME):$(TAG) --platform linux/amd64 .

push-image:
	@echo "Pushing Docker image $(IMAGE_NAME):$(TAG)..."
	docker push $(IMAGE_NAME):$(TAG)

run-image:
	@echo "Running Docker image $(IMAGE_NAME):$(TAG)..."
	docker run -p 3000:3000 $(IMAGE_NAME):$(TAG)

clean-image:
	@echo "Cleaning dangling Docker images..."
	docker rmi $(docker images --filter "dangling=true" -q) -f


# Kubernetes commands
.PHONY: helm helm-uninstall clean certs develop
develop: helm-uninstall build-image push-image helm

helm:
	@echo "Upgrading/Installing $(DEPLOYMENT) Helm chart..."
	helm upgrade --install $(DEPLOYMENT) ./helm --set image.tag=$(TAG) --namespace $(NAMESPACE)

helm-uninstall:
	@echo "Uninstalling $(DEPLOYMENT) from Kubernetes..."
	helm uninstall demo -n $(NAMESPACE) || true

clean:
	@echo "Cleaning up all resources in the $(NAMESPACE) namespace..."
	kubectl delete all --all -n $(NAMESPACE) || true
	kubectl delete namespace $(NAMESPACE) || true
	sleep 2

# Release and versioning
.PHONY: release bump-version update-changelog
release: update-changelog bump-version build-push

update-changelog:
	@echo "Updating changelog..."
	./scripts/update-changelog.sh

bump-version:
	@echo "Bumping version..."
	./scripts/bump-version.sh

build-push:
	@echo "Building and pushing Docker image $(IMAGE_NAME):$(VERSION)..."
	docker build -t $(IMAGE_NAME):$(VERSION) -t $(IMAGE_NAME):latest .
	docker push $(IMAGE_NAME):$(VERSION)
	docker push $(IMAGE_NAME):latest

# GitHub Pages and Helm chart publishing
.PHONY: gh-pages-publish helm-repo-update

gh-pages-publish:
	@echo "Publishing Helm chart for $(SERVICE_NAME) to GitHub Pages..."
	helm package ./$(CHART_DIRECTORY) -d /tmp
	git checkout gh-pages || git checkout -b gh-pages
	cp /tmp/$(DEPLOYMENT)-*.tgz .
	if [ -f index.yaml ]; then \
	helm repo index . --url https://$(REMOTE_CHART_REPOSITORY)/$(SERVICE_NAME)/ --merge index.yaml; \
	else \
	helm repo index . --url https://$(REMOTE_CHART_REPOSITORY)/$(SERVICE_NAME)/; \
	fi
	git add .
	git commit -m "chore: update Helm chart to v$(VERSION)"
	git push origin gh-pages
	git checkout main

helm-repo-update:
	@echo "Adding and updating Helm repo for $(SERVICE_NAME)..."
	helm repo add $(SERVICE_NAME) https://$(REMOTE_CHART_REPOSITORY)/$(SERVICE_NAME)/
	helm repo update
	helm repo list