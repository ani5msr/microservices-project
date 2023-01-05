FROM golang:1.18 AS builder

WORKDIR /build

ADD ./go.mod  go.mod
ADD ./main.go main.go
ADD ./service service

# Update
RUN apt-get --allow-releaseinfo-change update && apt upgrade -y

# Fetch dependencies
RUN go mod download all

# Build image as a truly static Go binary
RUN CGO_ENABLED=0 GOOS=linux go build -o /post_service -a -tags netgo -ldflags '-s -w' .

FROM scratch
COPY --from=builder /link_service /app/postdservice
EXPOSE 8080
ENTRYPOINT ["/app/link_service"]
