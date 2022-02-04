# A backend for the Gigs industry

[![Go Report Card](https://goreportcard.com/badge/github.com/pvr1/gigs)](https://goreportcard.com/report/github.com/pvr1/gigs)
[![pvr1](https://circleci.com/gh/pvr1/gigs.svg?style=svg)](https://github.com/pvr1/gigs)
[![Total alerts](https://img.shields.io/lgtm/alerts/g/pvr1/gigs.svg?logo=lgtm&logoWidth=18)](https://lgtm.com/projects/g/pvr1/gigs/alerts/)

## Run the backend

Install into ~/go/github.com/pvr1/gigs

Run by 

```bash
go run .
```

Dockerize by

```bash
docker build .
```

Create mongodb in kubernetes (in namespace mongodb). Create gigs in default as of now...

```bash
kubectl ns mongodb
helm repo add bitnami https://charts.bitnami.com/bitnami
helm install my-release -n mongodb bitnami/mongodb

kubectl apply -f k8s_gigs.yaml
```


