FROM golang:1.23
WORKDIR /app
ENV GOPROXY=off GOSUMDB=off CGO_ENABLED=0
COPY go.mod go.sum ./
COPY vendor ./vendor
COPY cmd ./cmd
COPY internal ./internal
RUN go build -mod=vendor -o /app/bldghvac ./cmd/bldghvac
EXPOSE 8080
CMD ["/app/bldghvac", "-addr", ":8080", "-data", "/app/data"]
