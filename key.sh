#!/bin/bash

export AWS_ACCESS_KEY_ID=foo
export AWS_SECRET_ACCESS_KEY=foo

aws --endpoint-url=http://localhost:4572 s3api create-bucket --bucket test-bucket --region=us-east-1