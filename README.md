# ctf

Capture the flag client and server for Kubernetes

``` bash
mkdir -p ~/.kube
nano ~/.kube/config1

export KUBECONFIG=~/.kube/config1

kubectl cluster-info
```


``` bash
kubectl delete -f ./pod.yaml

kubectl create -f ./pod.yaml

```