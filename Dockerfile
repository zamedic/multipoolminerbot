FROM alpine:3.6
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
WORKDIR /app
# Now just add the binary
COPY multipoolminerbot /app/
ENTRYPOINT ["/app/multipoolminerbot"]