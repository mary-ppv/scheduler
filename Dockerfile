FROM golang:1.23.2 AS builder

WORKDIR /scheduler

COPY . .

RUN go mod tidy

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /main 

FROM ubuntu:latest

WORKDIR /scheduler

COPY --from=builder /main /scheduler/main

COPY --from=builder /scheduler/web /scheduler/web

COPY --from=builder /scheduler/scheduler.db /scheduler/scheduler.db

COPY --from=builder /scheduler/.env /scheduler/.env

ENV TODO_PORT=${TODO_PORT}
ENV TODO_DBFILE=${TODO_DBFILE}

EXPOSE ${TODO_PORT}

CMD ["/scheduler/main"]
