# Build container
FROM golang:1.14 AS build

WORKDIR /go/src
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o backend ./cmd/backend

# Running container
FROM scratch

# Copy SSL certificates from build container.
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

EXPOSE 9000/tcp

VOLUME [ "/config" ]

# Environment variables requires by application. Default values.
ENV FIREBASE_CONFIG "/config/firebase-config.json"

WORKDIR /
COPY --from=build /go/src/backend /
CMD [ "/backend" ]