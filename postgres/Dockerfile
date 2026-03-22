FROM postgres:16

RUN apt-get update && apt-get install -y openssl && rm -rf /var/lib/apt/lists/*

RUN mkdir -p /var/lib/postgresql/ssl && \
    openssl req -new -x509 -days 365 -nodes \
      -out /var/lib/postgresql/ssl/server.crt \
      -keyout /var/lib/postgresql/ssl/server.key \
      -subj "/C=RU/ST=Moscow/O=SecureAudit/CN=postgres" && \
    chmod 600 /var/lib/postgresql/ssl/server.key && \
    chown postgres:postgres /var/lib/postgresql/ssl/server.key /var/lib/postgresql/ssl/server.crt
