FROM golang:latest as builder

RUN mkdir /code
WORKDIR /code/

ADD . .
ENV GO111MODULE=on
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/app

FROM alpine

WORKDIR /bin/

COPY --from=builder /bin/app .

CMD exec /bin/app