FROM scratch

WORKDIR /app

# COPY go.mod go.sum /app/

COPY public /app/public

COPY k8s-api-client /usr/bin/k8s-api-client

EXPOSE 8080

ENTRYPOINT ["/usr/bin/k8s-api-client"]