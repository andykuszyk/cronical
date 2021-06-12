FROM golang:alpine AS build
WORKDIR /cronical
COPY . ./
RUN go build -o cronical ./cmd/cronical/main.go

FROM alpine
COPY --from=build /cronical/cronical /cronical
COPY ./internal/cronical/html /html/
CMD /cronical
