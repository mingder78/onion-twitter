FROM asia.gcr.io/gogetdb/base-onion:latest 
MAINTAINER Ming-der Wang <ming@log4analytics.com>
CMD ["/go/bin/onion","serve"]
EXPOSE 8080
