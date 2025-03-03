# ctf

Capture the flag client and server for Kubernetes

## Connect to K8s cluster

``` bash
mkdir -p ~/.kube
nano ~/.kube/config1

export KUBECONFIG=~/.kube/config1

kubectl cluster-info
```

## Create a pod

``` bash
kubectl delete -f ./pod.yaml

kubectl create -f ./pod.yaml

```

## Check the pod

``` bash
kubectl get pods -n <your_namespace>

kubectl logs <pod_name> -n <your_namespace>

kubectl describe pod <pod_name> -n <your_namespace>
```
