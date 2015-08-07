#!/bin/bash
cd target/linux_i386/
file=tatanka
bucket="tnolet-tatanka"
resource="/${bucket}/${file}"
dateValue=`date -u +"%Y%m%dT%H%M%SZ"`
stringToSign="PUT\n\n${contentType}\n${dateValue}\n${resource}"
s3Key=$AWS_ACCESS_KEY_ID
s3Secret=$AWS_SECRET_ACCESS_KEY
signature=`echo -en ${stringToSign} | openssl sha1 -hmac ${s3Secret} -binary | base64`
curl -X PUT -T "${file}" \
  -H "Host: ${bucket}.s3.amazonaws.com" \
  -H "Date: ${dateValue}" \
  -H "Content-Type: ${contentType}" \
  -H "Authorization: AWS ${s3Key}:${signature}" \
  https://${bucket}.s3.amazonaws.com/${file}

    # -H "x-amz-acl: public-read" \
