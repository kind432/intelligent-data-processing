FROM golang:1.20 AS builder

WORKDIR /intelligent_data_processing

COPY . .

RUN go mod download

RUN go build -o app ./main.go

CMD ["/intelligent_data_processing/app"]
