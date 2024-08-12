#!/bin/bash

# 创建 cert 目录，保存 CA 文件
mkdir cert
cd cert

# 创建 CA 证书
openssl genrsa -out ca.key 2048
openssl req -x509 -new -nodes -key ca.key -subj "/CN=superproj" -days 3650 -out ca.crt
openssl genrsa -out server.key 2048
openssl req -new -key server.key -subj "/CN=superproj.com" -out server.csr
echo "subjectAltName = IP:10.37.83.200" > extfile.cnf
openssl x509 -req -in server.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out server.crt -days 365 -extfile extfile.cnf

# 输出根证书的内容，并填充到 webhooks[].clientConfig.caBundle
cat ca.crt | base64 -w 0
