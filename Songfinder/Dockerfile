FROM golang:latest

RUN go version
ENV GOPATH=/

EXPOSE 8080

COPY ./ ./

RUN go build -o Songfinder ./cmd/app/main.go

CMD ["./Songfinder"]