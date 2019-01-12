FROM ubuntu
COPY main .
COPY .env .
COPY key.pem .
COPY cert.pem .
COPY views views
COPY static static
COPY sessions tmp
EXPOSE 80
EXPOSE 443
ENTRYPOINT ["./main"]
