FROM golang:alpine as builder
WORKDIR /base
COPY . .
RUN go mod download
RUN go build -o quiz3

FROM alpine:latest
RUN mkdir /app
RUN mkdir /app/public

COPY --from=builder /base/quiz3 /app/quiz3
COPY prod.config.ini /app/conf.ini
ADD public /app/public/

WORKDIR /app

CMD ["./quiz3", "-m"]

