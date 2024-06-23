FROM golang:1.22

WORKDIR /go/src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /go/bin/main .

EXPOSE 8080

CMD ["/go/bin/main"]
