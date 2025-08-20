FROM golang:1.23 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o fileServer .

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/fileServer .

RUN mkdir -p /root/uploads

VOLUME [ "/root/uploads" ]

EXPOSE 8080

CMD [ "./fileServer" ]