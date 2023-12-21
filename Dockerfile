FROM golang:1.21.5-alpine3.19 AS build

WORKDIR /app
COPY . .

RUN go mod download
RUN go mod tidy
RUN go mod vendor
RUN go build -o /app/hello

FROM alpine:latest

WORKDIR /app
COPY --from=build /app/hello .

ENV APP_SECRET=kj3mSJbsw4lpFWUsHasQZf9r
ENV JWT_SECRET=L9C98ouj2SXUyRcz4HRn2sBwIIY5trlzIOyVkcBntWETBz7e4kbIYZwAuVyIBNkyw
CMD ["./hello"]