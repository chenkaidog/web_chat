# FROM ubuntu:latest

# WORKDIR /app

# COPY output/ .

# ENV log_level=info\
#     log_output_filename='web_chat'

# RUN chmod +x bootstrap.sh

# CMD ["./bootstrap.sh"]

FROM golang:1.21

WORKDIR /app

COPY . .

RUN mkdir /log \
    && go env -w GO111MODULE=on \
    && go env -w GOPROXY=https://goproxy.cn,direct \
    && go mod tidy \
    && go build -o main

ENV log_level=info\
    log_output_filename='web_chat'

CMD ["./main"]