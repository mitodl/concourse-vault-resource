FROM golang:1.26-alpine@sha256:2389ebfa5b7f43eeafbd6be0c3700cc46690ef842ad962f6c5bd6be49ed82039 as build
WORKDIR /go/src/github.com/mitodl/concourse-vault-plugin
COPY . .
RUN apk add make && make release

FROM alpine:3.23@sha256:25109184c71bdad752c8312a8623239686a9a2071e8825f20acb8f2198c3f659
WORKDIR /opt/resource
COPY --from=build /go/src/github.com/mitodl/concourse-vault-plugin/check .
COPY --from=build /go/src/github.com/mitodl/concourse-vault-plugin/in .
COPY --from=build /go/src/github.com/mitodl/concourse-vault-plugin/out .
