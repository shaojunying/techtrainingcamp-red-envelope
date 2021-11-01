FROM centos:7
COPY main /root/server
EXPOSE 8080
CMD /root/server