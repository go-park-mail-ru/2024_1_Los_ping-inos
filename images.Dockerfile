FROM golang:1.21.0-alpine AS build-stage

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download && go mod tidy
RUN go clean --modcache

COPY ./ ./

WORKDIR /app/internal/

RUN CGO_ENABLED=0 GOOS=linux go build -o /images ./cmd/images/main.go

# Deploy the application binary into a lean image
FROM gcr.io/distroless/base-debian11 AS build-release-stage

WORKDIR /

COPY --from=build-stage images images
COPY . .

EXPOSE 8082

USER nonroot:nonroot

ENTRYPOINT ["/images"]