FROM golang:1.21-alpine

WORKDIR /usr/src/app/

COPY . .

EXPOSE 8080

CMD ["go", "run", "/usr/src/app/server.go"]

