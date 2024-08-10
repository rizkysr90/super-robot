# We need busybox for its sh and tee for entrypoint and stdout logs redirect
FROM busybox:1.35.0-uclibc AS busybox
# debian12 = bookworm
FROM golang:1.22-bookworm AS builder
WORKDIR /go/src/builder
COPY . .
RUN make build/grpc

# https://github.com/GoogleContainerTools/distroless
FROM gcr.io/distroless/static-debian12:nonroot
# Copy busybox sh and tee required by entrypoint.sh
COPY --from=busybox /bin/sh /bin/sh
COPY --from=busybox /bin/tee /bin/tee
# Copy the server binary
COPY --from=builder /go/src/builder/build/grpc /app/server
COPY docker-entrypoint.sh /app/entrypoint.sh
WORKDIR /app
CMD ["/bin/sh", "./entrypoint.sh"]