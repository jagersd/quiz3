FROM alpine:latest
RUN mkdir /app
RUN mkdir /app/public

COPY quiz3 /app
COPY prod.config.ini /app/conf.ini
ADD public /app/public/

WORKDIR /app

CMD ["./quiz3", "-m"]

