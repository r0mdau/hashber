# hashber

My simple implementation of a distributed system with consistent hashing and a gossip protocol.

## Getting started

Start your minikube cluster and share docker context

```
minikube start
eval $(minikube -p minikube docker-env)
```

Build the docker image

```
make docker
```

Create the docker deployment

```
kubectl create deploy hashber --image=hashber:latest --replicas=1 --port=8090 --port=7946

# patch the deployment to not pull images and use minikube local images
kubectl patch deployment hashber --patch '{"spec": {"template": {"spec": {"containers": [{"name": "hashber","imagePullPolicy": "Never"}]}}}}'
```

Expose it via a service

```
kubectl expose deploy/hashber
```

Then look at the logs and change replicas number in another terminal

```
kubectl scale --replicas=3 deploy/hashber
```