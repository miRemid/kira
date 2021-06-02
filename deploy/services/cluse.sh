#!/bin/sh

docker-compose scale auth=2 user=2 file=2 upload=2 file-api=3 user-api=3 upload-api=3 site=3