FROM concourse/busyboxplus:base

# satisfy go crypto/x509
RUN cat /etc/ssl/certs/*.pem > /etc/ssl/certs/ca-certificates.crt

ADD backend /opt/app/backend

EXPOSE 4000

WORKDIR /opt/app/
CMD ["/opt/app/backend"]
