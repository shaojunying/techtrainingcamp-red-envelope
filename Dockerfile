FROM centos:7
COPY red_envelope /root/server
EXPOSE 8080
CMD /root/server