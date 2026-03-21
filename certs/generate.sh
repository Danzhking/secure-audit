#!/bin/bash
# Generate self-signed TLS certificate for Secure Audit Collector
openssl req -x509 -nodes -newkey rsa:2048 \
  -keyout server.key \
  -out server.crt \
  -days 365 \
  -subj "/C=RU/ST=Moscow/O=SecureAudit/CN=collector"
echo "TLS certificate generated: server.crt, server.key"
