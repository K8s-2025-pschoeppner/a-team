# Deploy an application to a Kubernetes cluster

## Create a pod from the image

### Create namespace

```bash
kubectl create namespace my-namespace
```

### Create a pod

Write a deployment configuration file, for example `deployment.yaml`:

```bash
kubectl create -f deployment.yaml

kubectl apply -f deployment.yaml
```

### Create LoadBalancer service

Write a service configuration file, for example `service.yaml`:

```bash
kubectl create -f service.yaml
````

Check with:

```bash
kubectl get svc -n a-team
```

Application should be available at `http://<EXTERNAL-IP>:<PORT>`
