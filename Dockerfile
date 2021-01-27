FROM iron/base
WORKDIR /app
# copy binary into image
COPY web/tmpl/ /app/web/tmpl/
COPY web/static/ /app/web/static/
COPY somfy-rts-gateway /app/
RUN ["chmod", "+x", "./somfy-rts-gateway"]
ENTRYPOINT ["./somfy-rts-gateway"]

EXPOSE 8080