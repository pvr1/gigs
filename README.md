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

```bash
kubectl get svc -n mongodb # look for or make LoadBalancer svc for mongodb
export MONGODB_ROOT_PASSWORD=$(kubectl get secret --namespace mongodb my-release-mongodb -o jsonpath="{.data.mongodb-root-password}" | base64 --decode)
mongosh admin --host 10.0.0.166 --authenticationDatabase admin -u root -p $MONGODB_ROOT_PASSWORD
use gigs
db.createUser({user: "gigbe",pwd:  "gigbe", roles: [ { role: "readWrite", db: "gigs" }]})
```

### Create zip file - to return files from gigs
package main

import (
    "archive/zip"
    "fmt"
    "io"
    "os"
)

func main() {
    fmt.Println("creating zip archive...")
    archive, err := os.Create("archive.zip")
    if err != nil {
        panic(err)
    }
    defer archive.Close()
    zipWriter := zip.NewWriter(archive)

    fmt.Println("opening first file...")
    f1, err := os.Open("test.csv")
    if err != nil {
        panic(err)
    }
    defer f1.Close()

    fmt.Println("writing first file to archive...")
    w1, err := zipWriter.Create("csv/test.csv")
    if err != nil {
        panic(err)
    }
    if _, err := io.Copy(w1, f1); err != nil {
        panic(err)
    }

    fmt.Println("opening second file")
    f2, err := os.Open("test.txt")
    if err != nil {
        panic(err)
    }
    defer f2.Close()

    fmt.Println("writing second file to archive...")
    w2, err := zipWriter.Create("txt/test.txt")
    if err != nil {
        panic(err)
    }
    if _, err := io.Copy(w2, f2); err != nil {
        panic(err)
    }
    fmt.Println("closing zip archive...")
    zipWriter.Close()
}