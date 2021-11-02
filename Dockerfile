FROM centos:7
COPY config /root/app/config
COPY database /root/app/database
COPY app.yml /root/app/
COPY main /root/app/
EXPOSE 8080
CMD /root/app/main