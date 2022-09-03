FROM golang:1.18

WORKDIR /app

COPY go.mod .
COPY go.sum .

COPY . .

RUN go mod download

RUN go build --tags dynamic -o benny .
CMD ["/app/benny"]