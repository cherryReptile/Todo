FROM golang:1.18

WORKDIR /usr/src/app

RUN go install github.com/cosmtrek/air@latest
RUN go install -tags 'postgres,sqlite3' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
RUN go install github.com/go-delve/delve/cmd/dlv@latest
RUN go install github.com/swaggo/swag/cmd/swag@latest