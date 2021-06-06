FROM scratch
COPY go-api-server /usr/bin/go-api-server
ENTRYPOINT ["/usr/bin/go-api-server"]
