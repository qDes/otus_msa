# HW2
simple app with one http handler in k8s <br>


## minikube macOS arm hacks
```
minikube start

minikube addons enable ingress

minikube tunnel
```

## apply manifests

```
cd manifests
kubectl apply -f .
```

also add arch.homework to /etc/hosts for better experience

## check deployment
```
curl http://arch.homework/otusapp/andrei/health 
```