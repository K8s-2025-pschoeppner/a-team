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

Write a service configuration file, for example `loadbalancer.yaml`:

```bash
kubectl create -f loadbalancer.yaml
````

Check with:

```bash
kubectl get svc -n a-team
```

Application should be available at `http://<EXTERNAL-IP>:<PORT>`

### Availability of the pod

Requirements:

- Multiple replicas of the pod
- More than one pod
- Pods on more than one node

#### Create PodDisruptionBudget (Not required but nice to make sure than a minimum number of pods are running)

The PodDisruptionBudget (PDB) is a policy that limits the number of disruptions that can take place simultaneously.
So that means if we configure it to run 3 Pods, it will ensure that at least 3 Pods are running at all times.

Write a PodDisruptionBudget configuration file, for example `pdb.yaml`:

```bash
kubectl create -f poddisruption.yaml
```

#### Config the loadbalancing

Configuring how pods are getting splitted over multiple Nodes (see `deployment.yaml`).

```bash
spec:
  replicas: 5
  selector:
    matchLabels:
      app: github-application
  template:
    metadata:
      labels:
        app: github-application
    spec:
      # Entferne jegliche requiredDuringSchedulingIgnoredDuringExecution Anti-Affinity
      topologySpreadConstraints:
      - maxSkew: 1
        topologyKey: kubernetes.io/hostname
        whenUnsatisfiable: ScheduleAnyway
        labelSelector:
          matchLabels:
            app: github-application
      affinity:
        podAntiAffinity:
          # Nur preferred verwenden, nicht required
          preferredDuringSchedulingIgnoredDuringExecution:
          - weight: 100
            podAffinityTerm:
              labelSelector:
                matchExpressions:
                - key: app
                  operator: In
                  values:
                  - github-application
              topologyKey: "kubernetes.io/hostname"
```
