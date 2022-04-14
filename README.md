# UPLOAD EXAMPLE
### Spec:
    - Create an http REST API server in GO for uploading data and then upload the following payload into a cloud storage such as s3 or google cloud storage.
    - The payload size limit is 10KB
    - API can still handle/receive requests from clients without having to wait for the uploading process to complete

### To run the test
```
make test
```
### To run the benchmark
```
make benchmark
```
The benchmark will take around 9s to run 100 request.

```
100	  99007572 ns/op	   39115 B/op	     468 allocs/op
PASS
ok  	github.com/eneoti/upload-example/apis	9.915s
```

The step upload file to S3/GCS will be assumed cost 1s.
```
cloudStorageClient.On("Upload",..).After(1 * time.Second).Return(nil)
```


### To run the server
```
make run
```
### To run the container
```
docker-compose up
```

### TODO:
1. The S3Hanlder and GCSHanlder are implemented as generic. It's still not test with AWS/GOOGLE.
Need to modify them.
2. This code is very simple design without using any framework. We should apply the framework out there later.
3. This code skips the refining the payload. We may add database interaction later.
4. Need to add authentication of the uploading data API.
5. Need to add environment configuation and apply it in all modules.

### Log
```
Use the https://github.com/uber-go/zap
```
### Explaination
1. Apply the RateLimit when uploading to S3,GCS:
	This code set the ratelimit = 10 as default. Please reference the [this link](https://aws.amazon.com/premiumsupport/knowledge-center/s3-503-within-request-rate-prefix/) to have the ratelimit information when working with AWS.

