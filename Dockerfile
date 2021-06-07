FROM scratch
COPY k8s-api-client /usr/bin/k8s-api-client
ENTRYPOINT ["/usr/bin/k8s-api-client"]
