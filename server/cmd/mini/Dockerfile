FROM golang:alpine AS builder
RUN apk --no-cache update && apk --no-cache add gcc musl-dev git make bash
WORKDIR /source
COPY ./server .
RUN GOOS=linux GOARCH=amd64 make mini
FROM alpine
RUN apk --no-cache update && apk --no-cache add ca-certificates
RUN addgroup -S app && adduser -S app -G app
COPY ./client/build /public
WORKDIR /bin
COPY --from=builder /source/bin/mini mini
RUN chown -R app:app /bin /public
ENV STATIC_FILES_DIR=/public
USER app
EXPOSE 9000
CMD ["./mini"]
