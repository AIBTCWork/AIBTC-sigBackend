FROM golang:1.25 AS builder

ENV GO111MODULE=on \
    TZ=Asia/Shanghai \
    GOPROXY=https://goproxy.cn,direct
    
ARG BUILD_MODE

# set working directory
WORKDIR /app

# copy go modules and dependencies
COPY ./go.mod ./go.sum ./
RUN go mod download

# copy project files
COPY . .

RUN if [ "${BUILD_MODE}" = "debug" ]; then \
        cp ./conf/config_debug.yaml ./config.yaml; \
    else \
        cp ./conf/config_prod.yaml ./config.yaml; \
    fi

# RUN  cp ./config_debug.yaml ./config.yaml
# build binary file
# -ldflags "-s -w"  去除编译信息 减小体积
# 注意：输出文件名使用 /app-server 避免与 server/ 目录冲突
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app-server -v -ldflags "-s -w" ./server

# use Alpine 3.19 as runtime environment
FROM alpine:3.19

# install timezone data and set timezone to Shanghai
RUN apk add --no-cache tzdata \
    && ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/Shanghai" > /etc/timezone

COPY --from=builder /app-server /server
COPY --from=builder /app/config.yaml /config.yaml
# give execution permission
RUN chmod +x /server

# CMD ["/server"]
# 最佳实践
ENTRYPOINT [ "/server" ]