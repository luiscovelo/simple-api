When using kind to creating k8s infrastructure locally.

It's necessary to expose manually the `nodePort` from cluster config because kind uses docker instead of VM.

```yaml
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
nodes:
- role: control-plane
  extraPortMappings:
  - containerPort: 30080
    hostPort: 180
    protocol: TCP
  - containerPort: 30443
    hostPort: 1443
    protocol: TCP
- role: worker
```

To use `ingress-nginx` it's necessary to set nodePorts when to installing it.

```bash
helm install ingress-nginx ingress-nginx/ingress-nginx \
  --namespace ingress-nginx \
  --create-namespace \
  --set controller.service.type=NodePort \
  --set controller.service.nodePorts.http=30080 \
  --set controller.service.nodePorts.https=30443 \
  --set controller.watchIngressWithoutClass=true \
  --wait
```

```bash
helm uninstall ingress-nginx -n ingress-nginx
```

To access dashboard, we need to create a forward port:

```bash
eyJhbGciOiJSUzI1NiIsImtpZCI6IlkwcFhES0hvTC1JQWtNc3NNYWN2cHRfaEtmdGVqWFZEZlRmclp4VXIzSkkifQ.eyJhdWQiOlsiaHR0cHM6Ly9rdWJlcm5ldGVzLmRlZmF1bHQuc3ZjLmNsdXN0ZXIubG9jYWwiXSwiZXhwIjoxNzQ0MTQxNTMzLCJpYXQiOjE3NDQxMzc5MzMsImlzcyI6Imh0dHBzOi8va3ViZXJuZXRlcy5kZWZhdWx0LnN2Yy5jbHVzdGVyLmxvY2FsIiwianRpIjoiMmEzY2QzMGItYzM5Zi00MzExLWE4YmMtY2EzNzNkZTUyYjA2Iiwia3ViZXJuZXRlcy5pbyI6eyJuYW1lc3BhY2UiOiJrdWJlcm5ldGVzLWRhc2hib2FyZCIsInNlcnZpY2VhY2NvdW50Ijp7Im5hbWUiOiJhZG1pbi11c2VyIiwidWlkIjoiNzVlMmU5MjEtMTk5NS00ZDY1LWIwMjEtZWM0ODkyZGZjZmMwIn19LCJuYmYiOjE3NDQxMzc5MzMsInN1YiI6InN5c3RlbTpzZXJ2aWNlYWNjb3VudDprdWJlcm5ldGVzLWRhc2hib2FyZDphZG1pbi11c2VyIn0.y6XHajCVFUSdcMrz2yO8xUW6i3Jc_470rD5GXfliQK8qXebw1d3rYD0gQfqgFjkDj_q3vt3ajJMyjkWdGWt2PwRPNTbe8aqKHL9CgGWVW9ozA6FfJQtDu7lWK0RmwY20CYs9elKuck3RO6xDs0LY2Lsw1AUV29uKdu6ahoxRDYgIcofuTSsih7MKJ76ryAuQ4mRRjxDebUpsh7N2CVMYsdc4E29Iyg49F_YbuLdIG_wo37bY9CKKSR6vKJWK1XRfgBgM6fBSdvftHUf6aJXdEjh9WDwef5lDmQFmcEhUHbW-fVKpB0022wLYxHi3mrU--8leMsjUDSHPZiBl1E3ssQ
```

```bash
kubectl -n kubernetes-dashboard port-forward svc/kubernetes-dashboard-kong-proxy 8443:443
```