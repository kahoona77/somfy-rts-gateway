FROM iron/base
WORKDIR /app
COPY somfy-rts-gateway /app/
RUN ["chmod", "+x", "./somfy-rts-gateway"]
ENTRYPOINT ["./somfy-rts-gateway"]

EXPOSE 8080
