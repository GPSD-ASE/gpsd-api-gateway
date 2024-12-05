## Network Flow Summary

1.	Client: Request to 127.0.0.1:54987 (Minikube tunnel with Nginx service on a dynamically created port).
   
2.	Minikube Tunnel: Maps a local port (54987) to the Kubernetes NodePort Nginx service.
   
3.	Nginx Pod:
   
	•	The request is routed to the Nginx pod through the NodePort service.

	•	Nginx acts as a reverse proxy and forwards the request to the API Gateway service.

4.	API Gateway Service:
   
	•	Nginx forwards the request to the ClusterIP of the API Gateway service, which routes traffic to one of the API Gateway pods.

5.	API Gateway Pod:
   
	•	The application running inside the API Gateway pod processes the request and sends a response.

6.	Response:
    
	•	The response flows back from the API Gateway pod → ClusterIP → Nginx pod → Minikube tunnel → client.

### Network Pathway
```
Client (curl/browser)
    |
    v
127.0.0.1:54987 (Minikube Tunnel)
    |
    v
NodePort Nginx Service (30080 on Minikube node)
    |
    v
Nginx Pod
    |
    v
ClusterIP API Gateway Service (10.108.83.91)
    |
    v
API Gateway Pod (10.244.0.57:3000)
    |
    v
Application (API Gateway: "API Gateway is running!")
```
