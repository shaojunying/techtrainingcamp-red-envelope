FROM centos:7
COPY main /root/server
COPY app.yml /root/
EXPOSE 8080
CMD /root/server