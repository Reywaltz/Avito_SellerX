FROM golang:latest

WORKDIR /go/src/github.com/Reywaltz/avito_backend

COPY . .

RUN go mod download

RUN go build -o ./bin/ ./cmd/avito_api/main.go

CMD [ "./bin/main" ]