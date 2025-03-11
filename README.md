# ctf

Capture the flag client and server for Kubernetes

## Connect to K8s cluster

``` bash
mkdir -p ~/.kube
nano ~/.kube/config1

export KUBECONFIG=~/.kube/config1

kubectl cluster-info
```

## Create pods/services

``` bash
kubectl delete -f ./pod.yaml

kubectl create -f ./pod.yaml

```

## Apply changes to pods/services

``` bash
kubectl delete pod github-action-pod -n a-team
kubectl apply -f ./pod.yaml
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

You find all related files in `roles` directory.

1. Create a service account
2. Create a role
3. Create a rolebinding
4. Add to the pod

``` bash
spec:
  serviceAccountName: ctf-serviceaccount
```

## Deploy an application with alreaady running ingress controller

1. Create a deployment
2. Create a service
3. Create an ingress

## Create different workload resources

Find all related files in `workload` directory.

- Deployment: Manages stateless applications by running a set of identical pods with automatic scaling, updates, and rollbacks.
- DaemonSet: Ensures a single pod runs on each (or selected) node, typically used for logging, monitoring, or node-specific services.
- CronJob: Executes pods periodically on a defined schedule, ideal for recurring tasks like backups or scheduled data processing.
- Job: Runs pods to completion for one-time, batch-processing tasks, such as database migrations or batch computations.

## Create a job disruption

Find all related files in `disruption` directory.

This ensure a minimun of X pods are running at all times.
In the PodDisruptionBudget you configure the minimum number of pods and the application this applies to.
