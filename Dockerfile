FROM golang:1.18-alpine AS build_base

RUN apk add --no-cache git  git gcc g++
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .

# ------------------------------------------------------------
# smaller image

RUN go build -o ./out/app

FROM alpine:latest
COPY --from=build_base /app/out/app /bin/restapi
EXPOSE 9090
CMD ["/bin/restapi"]