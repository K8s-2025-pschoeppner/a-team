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

## Restart pod

``` bash
kubectl delete pod github-action-pod -n a-team
kubectl apply -f ./kube/pod.yaml
```

## Check the cluster

``` bash
kubectl get pods -n <your_namespace>

kubectl get svc -n <your_namespace>

kubectl get ing -n <your_namespace>

kubectl logs <pod_name> -n <your_namespace>

kubectl describe pod <pod_name> -n <your_namespace>
```

## Add a serviceAccount to K8s

``` bash
kubectl apply -f ./kube/serviceaccount.yaml -n a-team

kubectl apply -f ./kube/role.yaml -n a-team

kubectl apply -f ./kube/rolebinding.yaml -n a-team
```

## Lösungswort

Kubernetes is often abbreviated to K8s because there are eight letters between the K and the S
