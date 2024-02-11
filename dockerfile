FROM ubuntu:latest


WORKDIR /app

COPY output/ .

RUN chmod +x bootstrap.sh

CMD ["./bootstrap.sh"]
