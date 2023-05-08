FROM golang:1.18-alpine
WORKDIR /opt/resource
COPY go.* /opt/resource
COPY cmd /opt/resource/
COPY pkg /opt/resource/
RUN go mod tidy && rm go.*
