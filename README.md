# k8s-api-client

The is a go client to  request for   the restart of the deployment resource in particular namespace

For slack Notification Set the variables

SLACK_WEBHOOK_URL
SLACK_CHANNEL
SLACK_USERNAME


run http server

```
make run
```

trigger restart for namespace(NS) and Deployment(Name)

```
curl -XPOST http://localhost:8080/restart -F NS=test -F Name=test-app

```
response 
```
[{"Name":"test-app","RestartedAt":"2021-07-10T14:19:13+07:00"}]
```

Make sure you have access to the  k8s api server


# to deploy in you cluster 

```
kustomize build kustomize/overlays | kubectl apply -f -
```