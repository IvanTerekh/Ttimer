FROM ubuntu
COPY main .
COPY .env .
COPY key.pem .
COPY cert.pem .
COPY veiws .
COPY static .
EXPOSE 8080
ENTRYPOINT ["./main"]
