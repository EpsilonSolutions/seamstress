FROM golang:1.16.7 as builder
ADD . /app
RUN cd /app && CGO_ENABLED=0 go build -mod vendor -o app

FROM scratch
COPY --from=builder /app/app /
CMD [ "/app" ]
