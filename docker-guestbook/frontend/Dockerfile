FROM concourse/busyboxplus:base

# satisfy go crypto/x509
RUN cat /etc/ssl/certs/*.pem > /etc/ssl/certs/ca-certificates.crt

ADD frontend /opt/app/frontend
ADD public /opt/app/public
ADD templates /opt/app/templates

EXPOSE 4000

WORKDIR /opt/app/
CMD ["/opt/app/frontend"]
