FROM golang:1.16.7 as builder
ADD . /app
RUN cd /app && go build -mod vendor -o app

FROM scratch
COPY --from=builder /app/app /
COPY entrypoint.sh entrypoint.sh
ENTRYPOINT [ "./app" ]
