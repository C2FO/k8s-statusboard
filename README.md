# k8s-statusboard

A dashboard to aggregate metrics and provide a quick view into the health of
multiple environments across different Kubernetes clusters & namespaces.

## Building & Running the docker image

`make image` will build a docker image called `k8s-statusboard`.

```
docker run -d \
           -p 8080:8080 \
           -v "$(pwd)/account.json:/root/.google/account.json" \
           -v "$(pwd)/config.yaml:/root/.kube/config" \
           -e "GOOGLE_APPLICATION_CREDENTIALS=/root/.google/account.json" \
           k8s-statusboard
```

This command will run the server on port 8080. This command expects that you
have a file in the current directory called account.json which contains your
google account credentials. It is also assumed that you have a file called
config.yaml which is your kubernetes config file.

## Notes

Tested using GKE clusters. If you are using a different credential provider,
it might be necessary to download the corresponding plugin and import it in
`pkg/k8s/config.go`.

