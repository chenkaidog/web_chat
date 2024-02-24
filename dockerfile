# FROM ubuntu:latest

# WORKDIR /app

# COPY output/ .

# RUN chmod +x bootstrap.sh

# CMD ["./bootstrap.sh"]

FROM golang:latest

WORKDIR /app

COPY . .

RUN mkdir /log \
    && go env -w GOPROXY='https://goproxy.io/' \
    && go mod tidy \
    && go build -o main

ENV log_level=info\
    log_output_filename='web_chat'

CMD ["./main"]