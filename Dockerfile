FROM alpine:latest

RUN apk update && apk add curl git

RUN curl -LO https://storage.googleapis.com/kubernetes-release/release/v1.21.1/bin/linux/amd64/kubectl
RUN chmod u+x kubectl && mv kubectl /bin/kubectl

WORKDIR /app


COPY k8s-api-client /usr/bin/k8s-api-client

EXPOSE 8080

ENTRYPOINT ["/usr/bin/k8s-api-client"]