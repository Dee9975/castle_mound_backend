FROM golang:1.16-buster AS build
WORKDIR /server
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY *.go .
RUN go build -o /castle_mound_server

FROM gcr.io/distroless/base-debian10
WORKDIR /
COPY --from=build /castle_mound_server /castle_mound_server
EXPOSE 8000
USER nonroot:nonroot
ENTRYPOINT ["/castle_mound_server"]