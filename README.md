# s3-dtr-tester
This tool is used to test Docker Trusted Registry's integration of S3 Functionality outside of the DTR UI itself.  It recreates the same test call made by Docker DTR within the `s3aws` `storageDriver` using the same version of AWS Go SDK that the `s3aws` `storageDriver` utilizes.


## Usage
Using the flags specify an AWS Access Key, Secret Key, Region, Endpoint URL and desired Bucket for testing.  A `PutObject` call will be made to create a file named with a UUID, the file will be checked for existence and if the test succeeds the file will be deleted.
```
docker run --rm -it squizzi/s3test -a <AWS Access Key> \
    -s <AWS Secret Key> \
    -r <AWS Region> \
    -e <AWS Endpoint URL> \
    -b <AWS Bucket>
```
