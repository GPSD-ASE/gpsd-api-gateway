## Network Flow Summary

1.	Client: Request to 127.0.0.1:54987 (Minikube tunnel).
2.	Minikube Tunnel: Maps local port (54987) to the Kubernetes NodePort service.
3.	NodePort Service: Routes traffic to the pod via its ClusterIP and TargetPort (3000).
4.	Pod: The application inside the pod processes the request and sends a response.
5.	Response: Flows back through the same path to the client.

### Network Pathway
```
Client (curl/browser)
    |
    v
127.0.0.1:54987 (Minikube Tunnel)
    |
    v
NodePort Service (30080 on Minikube node)
    |
    v
ClusterIP Service (10.108.83.91)
    |
    v
Pod (10.244.0.57:3000)
    |
    v
Application (API Gateway: "API Gateway is running!")
```