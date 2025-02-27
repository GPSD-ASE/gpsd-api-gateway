## GPSD API Gateway

The GPSD API Gateway serves as the central entry point for the disaster response system, handling authorization, and routing requests to appropriate microservices - gpsd-incident-mgmt, gpsd-map-mgmt and gpsd-user-mgmt.

### Table of Contents

	-	Features
	-	Architecture
	-	Project Structure
	-	Prerequisites
	-	Setting up Hashicorp Vault
	-	Environment Variables
	-	Running Locally
	-	Building and Running with Docker
	-	Deploying with Helm
	-	API Documentation
	-	Security
	-	License

### Features

	-	API Routing: Proxies requests to the appropriate backend services
	-	Security: TLS encryption with custom certificates
	-	Configuration: Environment-based configuration system
	-	Integration with HashiCorp Vault: Secure credential management

### Architecture

The API Gateway connects the following services:

-	User Management Service
-	Incident Management Service
-	Map Service
-	Notification Service (In progress)

### Project Structure

```bash
.
├── Dockerfile                 # Container definition
├── Makefile                   # Build automation
├── docs                       # Documentation
├── helm                       # Kubernetes Helm charts for deployment
├── internal                   # Go application code
│   ├── gateway                # Gateway specific code
│   │   └── pkg
│   │       ├── config         # Configuration management
│   │       ├── handlers       # Request handlers
│   │       └── routes         # API route definitions
│   └── main.go                # Application entry point
├── openapi                    # API specifications
└── private                    # Certificates and private configurations
```

### Prerequisites

-	Go 1.19+
-	Docker
-	Kubernetes cluster (for production deployment)
-	HashiCorp Vault

### Environment Variables

The API Gateway requires the following environment variables:

```bash
PORT=8080
USER_MGMT_HOST=user-mgmt-service
USER_MGMT_PORT=5500
INCIDENT_MGMT_HOST=incident-mgmt-service
INCIDENT_MGMT_PORT=9000
MAP_MGMT_HOST=map-mgmt-service
MAP_MGMT_PORT=8080
TLS_ENABLED=true
JWT_SECRET=your-jwt-secret
```
For Vault integration:
```bash
VAULT_ADDR=http://vault-ip:8200
VAULT_ROLE_ID=the-role-id
VAULT_SECRET_ID=the-secret-id
```

### Setting up Hashicorp Vault

The API Gateway can integrate with HashiCorp Vault for secure storage of sensitive information, particularly JWT secrets and X509 Certificates. Follow these steps to set up Vault and configure it for use with the API Gateway:


```bash
# Install Vault in Kubernetes using Helm
helm repo add hashicorp https://helm.releases.hashicorp.com
helm repo update
helm install vault hashicorp/vault --set "server.dev.enabled=true"

# Port-forward to access Vault
kubectl port-forward vault-0 8200:8200

# Set environment variables
export VAULT_ADDR='http://localhost:8200'
export VAULT_TOKEN='root'  # Only in dev mode, use proper authentication in production

# Enable Kubernetes auth
vault auth enable kubernetes

# Configure Kubernetes auth
vault write auth/kubernetes/config \
    kubernetes_host="https://$KUBERNETES_HOST:443" \
    token_reviewer_jwt="$(cat /var/run/secrets/kubernetes.io/serviceaccount/token)" \
    kubernetes_ca_cert="$(cat /var/run/secrets/kubernetes.io/serviceaccount/ca.crt)"

# Enable KV secrets engine v2
vault secrets enable -version=2 kv

# Create a JWT secret
JWT_SECRET="secret key"
vault kv put secret/gpsd/jwt secret_key="$JWT_SECRET"

# Create policy for API Gateway
cat > api-gateway-policy.hcl << EOF
path "secret/gpsd/jwt" {
  capabilities = ["read"]
}
EOF
vault policy write gpsd-api-gateway-policy api-gateway-policy.hcl

# Create a role for API Gateway
vault write auth/kubernetes/role/gpsd-api-gateway \
    bound_service_account_names=gpsd-api-gateway-sa \
    bound_service_account_namespaces=gpsd \
    policies=gpsd-api-gateway-policy \
    ttl=1h

# Create ServiceAccount in your Kubernetes manifest or apply:
kubectl apply -f - <<EOF
apiVersion: v1
kind: ServiceAccount
metadata:
  name: gpsd-api-gateway-sa
  namespace: gpsd
EOF

# Set up port-forwarding to API Gateway
kubectl port-forward gpsd/gpsd-api-gateway 3000:3000

# Make a request that uses JWT verification
curl -k https://localhost:3000/verify \
  -H "Authorization: Bearer [your-token]"
```

### Running Locally

##### Clone the repository

```bash
git clone https://github.com/gpsd-ase/gpsd-api-gateway.git
cd gpsd-api-gateway
```

##### Install dependencies

```bash
go mod tidy
go mod vendor
```

##### Run the application

```bash
go build -o gpsd-api-gateway ./internal/
bash gpsd-api-gateway
```

### Building and Running with Docker

##### Build the image

```bash
docker build -t gpsd/gpsd-api-gateway:latest .
```

##### Run the container

```bash
docker run -p 3000:3000 --env-file .env gpsd/gpsd-api-gateway:latest
```

### Deploying with Helm

##### Install Helm chart

```bash
helm install gpsd/gpsd-api-gateway ./helm
```

##### Upgrade existing deployment

```bash
helm upgrade gpsd/gpsd-api-gateway
```

### API Documentation
The API Gateway exposes the following endpoints:

##### User Management

-	POST /register - User registration
-	POST /register-admin - Admin registration
-	POST /login - User login
-	GET /verify - Token verification

##### Map Services

-	GET /zones - Get evacuation zones
-	GET /routing - Get route information
-	POST /evacuation - Create evacuation plan
-	GET /traffic - Get traffic information

##### Incident Management

-	GET /incidents - Get all incidents
-	POST /incidents - Create a new incident
-	GET /incidents/{id} - Get incident by ID
-	DELETE /incidents/{id} - Delete an incident
-	PATCH /incidents/{id}/status/{status} - Update incident status

### Security

The API Gateway implements several security measures:

-	JWT Authentication: All protected endpoints require a valid JWT token
-	TLS Encryption: HTTPS with custom certificates
-	HashiCorp Vault Integration: Secure storage for sensitive credentials
-	Role-Based Access Control: Different endpoints require different user roles

### License

This project is licensed under the MIT License - see the LICENSE file for details.