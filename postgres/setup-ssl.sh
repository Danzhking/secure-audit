#!/bin/bash
set -e

DATADIR="/var/lib/postgresql/data"

if [ ! -f "$DATADIR/server.key" ]; then
  openssl req -new -x509 -days 365 -nodes \
    -out "$DATADIR/server.crt" \
    -keyout "$DATADIR/server.key" \
    -subj "/C=RU/ST=Moscow/O=SecureAudit/CN=postgres"
  chmod 600 "$DATADIR/server.key"
  chown postgres:postgres "$DATADIR/server.key" "$DATADIR/server.crt"
fi
