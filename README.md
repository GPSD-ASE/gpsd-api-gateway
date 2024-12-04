## GPSD API Gateway

This repository contains the codebase for the GPSD API Gateway, designed to handle disaster response requests and serve OpenAPI documentation through an NGINX reverse proxy with self-signed SSL certificates for secure communication.

### Table of Contents
	•	Features
	•	Project Structure
	•	Prerequisites
	•	Setup Instructions
	•	Deployment on Kubernetes
	•	Testing the Setup
	•	Troubleshooting

### Features
•	Implements an API Gateway using Express.js.

•	Serves OpenAPI documentation via NGINX.

•	Supports secure communication with self-signed SSL certificates for local development.

•	Deployed on Kubernetes using Minikube for local testing.

### Project Structure
```
├── src/
│   └── app.js                # Express.js API Gateway server
├── openapi/
│   └── openapi.yaml          # OpenAPI Specification
├── deployments/
│   ├── api-gateway-deployment.yaml  # Kubernetes Deployment for API Gateway
│   ├── nginx-deployment.yaml        # Kubernetes Deployment for NGINX
├── services/
│   ├── api-gateway-service.yaml     # Kubernetes Service for API Gateway
│   ├── nginx-service.yaml           # Kubernetes Service for NGINX
├── docker/
│   ├── Dockerfile.api         # Dockerfile for the API Gateway
│   ├── Dockerfile.nginx       # Dockerfile for NGINX
├── utils/
│   └── startup_nginx.sh       # Script to initialize NGINX with SSL
├── README.md                  # Documentation
└── Makefile                   # Automation for setup and clean-up
```

### Prerequisites

1.	Docker: Ensure Docker is installed and running.

2.	Minikube: A local Kubernetes cluster.

3.	kubectl: Kubernetes command-line tool.

### Setup Instructions

1. Clone the Repository
```
git clone https://github.com/<your-username>/gpsd-api-gateway.git
cd gpsd-api-gateway
```

2. Build Docker Images

Build the images for the API Gateway and NGINX.
```
make build-image
```

3. Push Images (Optional)

If you’re using a container registry, push the images:
```
make push-image
```

### Deployment on Kubernetes

1. Start Minikube
```
minikube start --driver=docker
```

2. Deploy the Application

Use the Makefile to deploy all resources:
```
make all
```

This will:

•	Deploy the API Gateway and NGINX pods.

•	Set up Kubernetes services for both components.

•	Create self-signed SSL certificates during the NGINX container startup.

### Testing the Setup

1. Access OpenAPI Specification
```
curl https://api.gpsd.com/openapi.yaml -k
```

2. Test the API Gateway
```
curl https://api.gpsd.com/api/ -k
```

If everything is configured correctly, you should receive:
```
API Gateway is running!
```

### Troubleshooting

1.	Minikube IP Issues
  
Ensure the domain api.gpsd.com resolves to your Minikube IP:
```
echo "$(minikube ip) api.gpsd.com" | sudo tee -a /etc/hosts
```

2.	Certificate Issues

If certificates fail, verify NGINX logs:
```
kubectl logs -l app=gpsd-nginx -n gpsd
```

3.	Kubernetes Resources Not Found

Clean up and redeploy using:
```
make clean
make all
```

License

This project is yet to be licensed.
