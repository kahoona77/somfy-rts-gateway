FROM iron/base
WORKDIR /app
# copy binary into image
COPY web/tmpl/ /app/web/tmpl/
COPY web/static/ /app/web/static/
COPY somfy-rts-gateway /app/

ENTRYPOINT ["./somfy-rts-gateway"]

EXPOSE 8080