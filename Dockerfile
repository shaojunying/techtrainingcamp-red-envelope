FROM centos:7
WORKDIR /root/app
COPY red_envelope server
COPY app.yml ./
EXPOSE 8080
CMD ./server