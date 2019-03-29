FROM alpine:3.9
WORKDIR /function/
COPY bin/oraclecloud-billing-slack .
RUN apk update && \
	apk add --no-cache ca-certificates
ENTRYPOINT ["/function/oraclecloud-billing-slack"]
