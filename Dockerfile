# Build container
FROM golang:1.14 AS build

WORKDIR /go/src
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o backend ./cmd/backend

# Running container
FROM scratch

EXPOSE 9000/tcp

VOLUME [ "/config" ]

# Environment variables requires by application. Default values.
ENV GOOGLE_APPLICATION_CREDENTIALS "/config/google-services.json"
ENV FIREBASE_CONFIG "/config/firebase-config.json"

WORKDIR /
COPY --from=build /go/src/backend /
CMD [ "/backend" ]