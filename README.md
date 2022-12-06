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

Expose memberlist via a service

```
kubectl expose deploy/hashber --name=hashber-memberlist --port=7946
```

Expose hashber http API via a NodePort service

```
kubectl expose deploy/hashber --name=hashber --port=8090 --type=NodePort
```

Then look at the logs and change replicas number in another terminal

```
kubectl scale --replicas=3 deploy/hashber
```

Contact the hello API, you need to find your minikube master ip

```
# get nodeport port
kubectl get svc

# call aPI
curl http://192.168.49.2:31290/hello

# launch multiple calls
for i in {1..1000}; do sleep 1;curl http://192.168.49.2:31290/hello; done
```

Look at the logs to see who receive the request and who respond to it.

How to re deploy a new version

```
make docker
kubectl patch deployment $1 -p "{\"spec\": {\"template\": {\"metadata\": { \"labels\": {  \"redeploy\": \"$(date +%s)\"}}}}}"
```