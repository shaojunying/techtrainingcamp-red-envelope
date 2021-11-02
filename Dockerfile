FROM centos:7
RUN mkdir /root/config
COPY config /root/config/
RUN mkdir /root/database
COPY database /root/database/
COPY main /root/main
COPY app.yml /root/
EXPOSE 8080
CMD /root/main