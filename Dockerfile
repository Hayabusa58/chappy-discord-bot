# Build stage
FROM golang:1.23 AS build-stage

WORKDIR /usr/src/app

COPY .env ./
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY *.go ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /app
RUN go build -v -o /usr/local/bin/app ./...

# Prd stage
FROM gcr.io/distroless/base-debian11 AS prd-stage

WORKDIR /
COPY --from=build-stage /app /app

USER nonroot:nonroot
CMD ["/app"]