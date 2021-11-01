FROM centos:7
COPY config /root/
COPY database /root/
COPY main /root/main
COPY app.yml /root/
EXPOSE 8080
CMD /root/main