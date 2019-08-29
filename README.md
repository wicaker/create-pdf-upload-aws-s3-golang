# Create PDF from HTML and Upload to AWS S3 With Golang
This repositories was created to provide the way to create pdf from html template and upload to AWS S3 bucket. We will using `localstack` to get simulation of AWS S3 in local computer.


## How to run ?
Build and Run localstack:
- Just run '`docker-compose up`' in terminal.

Create credential AWS in local such as : AWS_ACCESS_KEY_ID , AWS_SECRET_ACCESS_KEY , and Bucket.
- Just run '`./key.sh`' in terminal.

Know we ready to create our pdf by running:
- '`go run main.go`'

To check our bucket :
- '`aws --endpoint-url=http://localhost:4572 s3 ls s3://`'

TO check our file in bucket :
- `aws --endpoint-url=http://localhost:4572 s3 ls s3://{our bucket name}`. For example : '`aws --endpoint-url=http://localhost:4572 s3 ls s3://test-bucket`'

Hope,  this repository may help you. Please contact me if something wrong with what i create.

GOOD LUCK :)
