FROM alpine

WORKDIR /app
# 需要先本地编译，手动 GOOS=linux GOARCH=amd64 go build -o openai
COPY openai .

RUN apk add tzdata && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/Shanghai" > /etc/timezone \
    && apk del tzdata

EXPOSE "$PORT"

CMD ["./openai"]
