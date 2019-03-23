FROM golang:1.12 as builder

WORKDIR /home/tinyrstats

COPY . .
RUN make skip_dep=false build



FROM alpine:3.8

WORKDIR /home/tinyrstats
COPY config.defaults.yml /etc/tinyrstats/config.yml
COPY sites.example.txt .
COPY scripts/entrypoint.sh /
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
RUN chmod 755 /entrypoint.sh && mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2 && adduser -D tinyrstats

USER tinyrstats

COPY --from=builder /home/tinyrstats/bin/tinyrstats .

EXPOSE 8080
ENTRYPOINT ["/entrypoint.sh", "/home/tinyrstats/tinyrstats"]
CMD ["monitor", "serve", "--preload-from-file", "./sites.example.txt"]
