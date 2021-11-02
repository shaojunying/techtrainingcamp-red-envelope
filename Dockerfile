FROM centos:7
COPY . /root/app
EXPOSE 8080
CMD /root/app/main