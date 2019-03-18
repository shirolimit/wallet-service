
FROM golang

WORKDIR /go/src/github.com/shirolimit/wallet-service

ADD . .

RUN go get -d -v ./...
RUN go install -v ./...

ENV CONNECTION_STRING ""

EXPOSE 8080/tcp

ENTRYPOINT wallet_service --connection-string=$CONNECTION_STRING --http-address=":8080"
