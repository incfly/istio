# The following section is used as base image if BASE_DISTRIBUTION=default
FROM ubuntu:bionic

COPY ts-server /usr/local/bin/

COPY server.crt /server.crt
COPY server.key /server.key

ENTRYPOINT ["/usr/local/bin/ts-server"]
