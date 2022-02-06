nerdctl build . -t pervonrosen/gigs:latest
nerdctl push pervonrosen/gigs:latest
kubectl get po --no-headers=true | awk '/gigs/{print $1}' | xargs  kubectl delete po

#kubectl get po --no-headers=true | awk '/^[[:blank:]]*gigs-/{print $1}' | xargs  kubectl logs 