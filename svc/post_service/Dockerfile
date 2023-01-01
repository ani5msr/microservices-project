FROM golang:1.16 AS builder
ADD ./main.go main.go
ADD ./service service
# Fetch dependencies
RUN go get -d -v

# Build image as a truly static Go binary
RUN CGO_ENABLED=0 GOOS=linux go build -o /post_service -a -tags netgo -ldflags '-s -w' .

FROM scratch
LABEL key="ani5msr"
COPY --from=builder /post_service /app/post_service
EXPOSE 7070
ENTRYPOINT ["/app/post_service"]
