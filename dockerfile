FROM golang as compiler
RUN  go get -a -ldflags '-s' \
github.com/Hleb112/krip_bot
FROM scratch
ENV LANGUAGE="en"
ENV TZ=Europe/Moscow
COPY --from=compiler /go/bin/krip_bot .
CMD ["/krip_bot"]