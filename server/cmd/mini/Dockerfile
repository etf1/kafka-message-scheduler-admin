FROM golang:alpine as builder
RUN apk --no-cache update && apk --no-cache add gcc musl-dev git make bash
WORKDIR /source
COPY ./server .
RUN ls -l
RUN GOOS=linux GOARCH=amd64 make mini

FROM alpine
RUN apk --no-cache update
COPY ./client/build /public
WORKDIR /bin
COPY --from=builder /source/bin/mini mini
ENV STATIC_FILES_DIR=/public
CMD ["./mini"]