NAMESPACE = gpsd
DEPLOYMENT = gpsd-api-gateway
IMAGE_NAME = $(NAMESPACE)/$(DEPLOYMENT)
TAG ?= latest  # If no tag is provided, default to 'latest'
SERVICE_NAME = $(DEPLOYMENT)

# Use `make develop` for local testing
develop: helm-uninstall build-image push-image helm

docker: build-image push-image

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
	git checkout gh-pages
	git rm -rf .
	git checkout HEAD -- helm
	git commit --allow-empty -m "fix: commit to update GitHub Pages"
	git push origin gh-pages -f
	cd helm/
	helm package .
	helm repo index . --url https://gpsd-ase.github.io/$(SERVICE_NAME)/ --merge index.yaml
	cd ..
	mv helm/$(SERVICE_NAME)-0.1.0.tgz helm/index.yaml .
	git rm -rf helm
	git add .
	git commit -m "fix: Publish Helm chart"
	git push origin gh-pages -f

helm-repo-update:
	@echo "Adding and updating Helm repo for $(SERVICE_NAME)..."
	helm repo add $(SERVICE_NAME) https://gpsd-ase.github.io/$(SERVICE_NAME)/
	helm repo update
	helm repo list