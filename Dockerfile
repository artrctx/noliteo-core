FROM golang:1.25-alpine AS builder
ARG CGO_ENABLED=0
WORKDIR /app

COPY go.mod go.sum ./
# DOWNLOAD DEPENDENCIES
RUN go mod download
COPY . .

# BUILD APP
RUN go build -o main cmd/api/main.go

FROM scratch
COPY --from=builder /app/main /main
ENTRYPOINT ["/main"]