#
# Copyright 2021 OpsMx, Inc.
#
# Licensed under the Apache License, Version 2.0 (the "License")
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#   http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#

#
# Install the latest versions of our mods.  This is done as a separate step
# so it will pull from an image cache if possible, unless there are changes.
#
FROM golang:1.17-alpine AS buildmod
ENV CGO_ENABLED=0
RUN mkdir /build
WORKDIR /build
COPY go.mod .
COPY go.sum .
RUN go mod download

#
# Compile the agent.
#
FROM buildmod AS build-binaries
COPY . .
RUN touch pkg/tunnel/tunnel.pb.go
RUN mkdir /out /out/agent-binaries
RUN go build -ldflags="-s -w" -o /out/agent app/agent/*.go
RUN go build -ldflags="-s -w" -o /out/controller app/controller/*.go
RUN go build -ldflags="-s -w" -o /out/make-ca app/make-ca/*.go
# Also build the different OS versions here, so we can publish them.
RUN GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o /out/agent-binaries/agent.amd64.latest app/agent/*.go
RUN GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -o /out/agent-binaries/agent.arm64.latest app/agent/*.go

#
# Establish a base OS image used by all the applications.
#
FROM alpine:3 AS base-image
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
RUN update-ca-certificates
RUN mkdir /local /local/ca-certificates && rm -rf /usr/local/share/ca-certificates && ln -s  /local/ca-certificates /usr/local/share/ca-certificates
COPY docker/run.sh /app/run.sh
ENTRYPOINT ["/bin/sh", "/app/run.sh"]

#
# For a base image without an OS, this can be used:
#
#FROM scratch AS base-image
#COPY --from=alpine:3 /etc/ssl/cert.pem /etc/ssl/cert.pem

#
# Build the agent image.  This should be a --target on docker build.
#
FROM base-image AS agent-image
WORKDIR /app
COPY --from=build-binaries /out/agent /app
EXPOSE 9102
CMD ["/app/agent"]

#
# Build the controller image.  This should be a --target on docker build.
# Note that the agent is also added, so the binary can be served from
# the controller to auto-update the remote agent.
#
FROM base-image AS controller-image
WORKDIR /app
COPY --from=build-binaries /out/controller /app
COPY --from=build-binaries /out/agent-binaries /app/agent-binaries
EXPOSE 9001-9002 9102
CMD ["/app/controller"]

#
# Build the make-ca image.  This should be a --target on docker build.
#
FROM base-image AS make-ca-image
WORKDIR /app
COPY --from=build-binaries /out/make-ca /app
CMD ["/app/make-ca"]
