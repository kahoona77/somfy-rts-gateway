FROM ubuntu:20.04
WORKDIR /app
# copy binary into image
COPY somfy-rts-gateway /app/

ENTRYPOINT ["./somfy-rts-gateway"]

EXPOSE 8080