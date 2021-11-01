FROM centos:7
ADD . /root/app
EXPOSE 8080
CMD /root/app/main