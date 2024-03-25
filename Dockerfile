FROM golang:1.21.3 AS builder

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o main .

CMD [ "./main" ]

EXPOSE 8089