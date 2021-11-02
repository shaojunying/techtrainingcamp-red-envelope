FROM centos:7
RUN mkdir -p /root/config
COPY config /root/config/
RUN mkdir -p /root/database
COPY database /root/database/
COPY main /root/main
COPY app.yml /root/
EXPOSE 8080
CMD /root/main