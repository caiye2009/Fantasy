FROM golang:1.25.4-alpine

WORKDIR /app
COPY . .
RUN apk add git build-base libc6-compat

#ENV GOPROXY=https://goproxy.cn,direct

RUN make build

EXPOSE 8080
CMD ["./fantasy-back"]
