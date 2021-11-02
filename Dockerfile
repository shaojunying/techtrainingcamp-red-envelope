FROM centos:7
RUN mkdir -p /root/config
ADD config /root/config
COPY database /root/
COPY main /root/main
COPY app.yml /root/
EXPOSE 8080
CMD /root/main