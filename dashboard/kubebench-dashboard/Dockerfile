FROM golang:onbuild
RUN mkdir /app/
ADD . /app/
WORKDIR /app
RUN go build -o main .
EXPOSE 9303
CMD ["/app/main"]