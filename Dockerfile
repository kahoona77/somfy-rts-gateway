FROM iron/base
WORKDIR /app
# copy binary into image
COPY web/tmpl/ /app/web/tmpl/
COPY web/assets/build/ /app/web/assets/build/
COPY emerald /app/

ENTRYPOINT ["./emerald"]

EXPOSE 8080