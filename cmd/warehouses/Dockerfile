FROM ubuntu:latest

WORKDIR /app

COPY /bin/warehouses/warehouses /app/warehouses/warehouses

COPY config/config.yaml /app

COPY /migrations /app/migrations

CMD ["/app/warehouses/warehouses"]