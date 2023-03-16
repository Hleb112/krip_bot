FROM golang:latest
ENV LANGUAGE="en"
ENV TZ=Europe/Moscow
RUN mkdir /app
ADD . /app/
WORKDIR /app
RUN go build -o main .
CMD ["/app/main"]