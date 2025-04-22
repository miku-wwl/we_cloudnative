FROM harbor.com/containers/golang:1.22.3-alpine3.20 as builder
WORKDIR /
COPY . .
RUN [ ! -d "vendor" ] && go env -w GOPROXY=https://goproxy.cn,direct || echo "skipping"
RUN go build -o bk

FROM harbor.com/containers/alpine:3.20-bookinfo-base
WORKDIR /
COPY templates /templates
COPY static /static
COPY cfg/config.yaml /cfg/config.yaml
COPY --from=builder /bk /
ENTRYPOINT [ "/bk", "serve", "--config", "cfg/config.yaml" ]