FROM alpine:3.6
WORKDIR /app
# Now just add the binary
COPY multipoolminerbot /app/
ENTRYPOINT ["/app/multipoolminerbot"]