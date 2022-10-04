FROM golang:latest as builder

WORKDIR /

COPY go.mod go.sum ./

RUN go mod download

COPY . .

WORKDIR /

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

### start new stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# COPY the pre-built binary
COPY --from=builder / .

EXPOSE 9010

CMD ["./main"]
