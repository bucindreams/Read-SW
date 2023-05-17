FROM golang:latest 

RUN mkdir /read-sw/
WORKDIR /read-sw/

COPY . .

RUN go build -o main main.go

CMD ["/go/main"]
