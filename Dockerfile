# Build container
FROM golang:1.14 AS build

WORKDIR /go/src
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o backend ./cmd/backend

# Running container
FROM scratch

EXPOSE 9000/tcp

# Environment variables requires by application. Default values.
ENV FIREBASE_CONFIG "/firebase-config.json"

WORKDIR /
COPY --from=build /go/src/backend /
COPY --from=build /go/src/firebase-config.json /
CMD [ "/backend" ]