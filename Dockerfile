FROM golang:alpine AS builder
WORKDIR /build
ADD go.mod .
COPY . .
ENV DOCKER_ENV=true
RUN go build -o msngr-backend cmd/apiserver/main.go
FROM alpine
WORKDIR /build
COPY --from=builder /build/msngr-backend /build/msngr-backend
COPY --from=builder /build/config /build/config
CMD ["./msngr-backend"]