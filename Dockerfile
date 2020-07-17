FROM alpine:latest

RUN mkdir /app

WORKDIR /app

ADD build/docker/odb-dicom /app/
