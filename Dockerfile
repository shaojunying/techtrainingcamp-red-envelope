FROM nginx
COPY * /root/
EXPOSE 8080
CMD /root/main