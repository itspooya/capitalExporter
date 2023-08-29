FROM golang:1.21 AS builder

LABEL authors="itspooya"
WORKDIR /app
COPY src/* /app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

FROM scratch
COPY --from=builder /app/main /app/
ENTRYPOINT ["/app/main"]
