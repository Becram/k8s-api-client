# k8s-api-client

This is aREST api to trigger the restart deployment resource in particular namespace

For slack Notification set
```
SLACK_WEBHOOK_URL
SLACK_CHANNEL
SLACK_USERNAME
```

run http server

```
make run
```

Trigger restart for namespace(NS) and Deployment(Name)

```
curl -XPOST http://localhost:8080/restart -F NS=test -F Name=test-app

```

Response 
```
[{"Name":"test-app","RestartedAt":"2021-07-10T14:19:13+07:00"}]
```

Make sure you have access to the  k8s api server


# To deploy in your cluster 

```
kustomize build kustomize/overlays | kubectl apply -f -
```