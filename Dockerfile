FROM golang:1.19-bullseye

COPY build/* /usr/local/bin/
RUN chmod 755 /usr/local/bin/*

ENTRYPOINT ["hashber"]