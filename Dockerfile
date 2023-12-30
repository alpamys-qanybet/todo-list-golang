FROM golang:1.19-bullseye AS golangbuild

WORKDIR /app
COPY . .

RUN go mod download
RUN go mod tidy
RUN go mod vendor
RUN go build cmd/app/main.go
RUN mv main /app/hello
RUN chmod +x /app/scripts/wait-for-postgres.sh

FROM debian:bullseye

# install psql
# why? because we need to wait until postgres starts and then start our app
# so we need psql command to execute postgres started status check
RUN apt-get update
RUN apt-get -y install postgresql-client

COPY --from=golangbuild /app/scripts/wait-for-postgres.sh .
COPY --from=golangbuild /app/hello .

ENV DATABASE_URL=postgresql://postgres:postgres@host.docker.internal:5433/todo
ENV JWT_SECRET=L9C98ouj2SXUyRcz4HRn2sBwIIY5trlzIOyVkcBntWETBz7e4kbIYZwAuVyIBNkyw
CMD ["./hello"]