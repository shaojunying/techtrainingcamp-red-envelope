FROM centos:7
COPY red_envelope /root/server
COPY app.yml /root/
EXPOSE 8080
CMD /root/server