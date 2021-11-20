FROM golang:1.13-alpine
ENV GOPROXY https://goproxy.cn,direct
ENV GO111MODULE on

WORKDIR /root/techtrainingcamp-red-envelop
COPY ./ ./
RUN go build -o ./server .
EXPOSE 8080
CMD ./server