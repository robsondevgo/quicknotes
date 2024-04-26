FROM golang:1.22-alpine
WORKDIR /app
RUN apk update && apk upgrade && apk add --no-cache ca-certificates
RUN update-ca-certificates
COPY . .
RUN go build -v -o ./server ./cmd/http/

FROM scratch
WORKDIR /bin
COPY --from=0 /app/server server
COPY --from=0 /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY .env .env
CMD ["/bin/server"]
