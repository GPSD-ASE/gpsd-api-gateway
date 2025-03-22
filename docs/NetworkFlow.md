### Network Flow Summary of GPSD Application
@author: Kaaviya Ramkumar

@email: prkaaviya17@gmail.com

@date: 22.03.2025


1. **Client**: Request to gpsd.duckdns.org (External domain pointing to ip 152.53.124.121)

2. **Traefik LoadBalancer**: 
   - Receives traffic on port 443 (HTTPS)
   - Handles TLS termination with Let's Encrypt certificates
   - Routes requests based on Host header

3. **API Gateway Service**:
   - Traefik forwards the request to the ClusterIP of the API Gateway service (gpsd-api-gateway)
   - Service routes traffic to one of the API Gateway pods

4. **API Gateway Pod**:
   - The application running inside the API Gateway pod processes the request
   - Communicates with other backend services (user-mgmt, map-mgmt, etc.)
   - Authenticates with Vault for secrets

5. **Response**:
   - The response flows back from the API Gateway pod → ClusterIP → Traefik LoadBalancer → Client

### Network Pathway

```
Client (curl/browser)
    |
    v
gpsd.duckdns.org (DNS points to 152.53.124.121)
    |
    v
Traefik LoadBalancer (152.53.124.121:443)
    |
    v
ClusterIP API Gateway Service (example - 10.43.33.214:80)
    |
    v
API Gateway Pod (example - 10.42.0.229:3000)
    |
    v
Application (API Gateway: "API Gateway is running at port 3000!")
    |
    v
Backend Services (User Management, Map Management, etc.)
    |
    v
Vault (Secret Management example - 10.43.176.90:8200)
```
