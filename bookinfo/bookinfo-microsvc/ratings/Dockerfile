FROM crpi-bbsyw9wy6j39ufia.cn-beijing.personal.cr.aliyuncs.com/containers-wwl/golang:1.22.3-alpine3.20 as builder
WORKDIR /
COPY . .
RUN [ ! -d "vendor" ] && go env -w GOPROXY=https://goproxy.cn,direct || echo "skipping"
RUN go build -o bk

FROM crpi-bbsyw9wy6j39ufia.cn-beijing.personal.cr.aliyuncs.com/containers-wwl/alpine:3.20-bookinfo-base
WORKDIR /
COPY templates /templates
COPY static /static
COPY cfg/config.yaml /cfg/config.yaml
COPY --from=builder /bk /
ENTRYPOINT [ "/bk", "serve", "--config", "cfg/config.yaml" ]