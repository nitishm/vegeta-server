FROM golang:1.11
ENV ROOT=/vegeta-server
ADD . $ROOT
RUN chmod 777 $ROOT/docker-entrypoint.sh
WORKDIR $ROOT
RUN make build
ENTRYPOINT ["./docker-entrypoint.sh"]
EXPOSE 80
