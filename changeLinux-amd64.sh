#!/bin/sh

sed -i 's/darwin/linux/g' openjpeg/j2kcodec.go
sed -i 's/darwin/linux/g' jpeglib/codec8.go
sed -i 's/darwin/linux/g' jpeglib/codec12.go
sed -i 's/darwin/linux/g' jpeglib/codec16.go

rm -rf jpeglib/dcmjpeg/linux
cp -a jpeglib/dcmjpeg/linux-amd64 jpeglib/dcmjpeg/linux

rm -rf openjpeg/j2klib/linux
cp -a openjpeg/j2klib/linux-amd64 openjpeg/j2klib/linux
