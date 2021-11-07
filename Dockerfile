FROM centos:7
WORKDIR /root/app
COPY red_envelope server
COPY config.yaml ./
EXPOSE 8080
CMD ./server