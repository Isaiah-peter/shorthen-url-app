FROM golang:alphine as builder

RUN mkdir /build

ADD . /build/

WORKDIR /build

RUN go build -o main

#stage 2

FROM alphine

RUN adduser -S -D -H /app appuser

USER appuser

COPY . /app

COPY --from=builder /build/main /app/

WORKDIR /app

EXPOSE 3000

CMD [ "./main" ]


