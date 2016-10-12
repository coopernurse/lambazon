
# Overview

This is an example Go web app that uses the [gin](https://godoc.org/github.com/gin-gonic/gin)
framework, which is a standard `net/http` web framework that knows nothing about AWS Lambda.

This app can be run and tested locally just like any other Go web app
and then deployed as a single AWS Lambda function.

This Lambda function can then be invoked using the Caddy web server as
a gateway using the [caddy-awslambda plugin](https://github.com/coopernurse/caddy-awslambda).

## Prerequisites

You need:

* A computer with Go installed
* An AWS account
* Docker (to simplify running Caddy)

## Run locally

During development, run the app locally. From this dir, run:

```
go get github.com/gin-gonic/gin
go get github.com/apex/go-apex
go get github.com/coopernurse/lambazon
go run calcweb/main/calcmain.go
```

The app should be available at: http://localhost:8080/

In a normal workflow you'd modify the app and stop/start the web server
without having to use Lambda.

## Deploy with Apex

[Apex]() is a command line tool that simplifies the process of registering and
deploying AWS Lambda functions.

* Install Apex per the instructions on the site
* Create an IAM role with the standard `AWSLambdaBasicExecutionRole` policy attached
* Edit the `project.json` and edit `role` to the ARN of the above IAM role
* Run: `apex deploy`
  * This will compile `functions/web/main.go` and deploy it as a `calc_web` Lambda function

## Run Caddy with `awslambda` plugin

Set the AWS env vars and invoke `run.sh` to start the Caddy web server.

```
export AWS_REGION=us-west-2
export AWS_ACCESS_KEY_ID=xyz
export AWS_SECRET_ACCESS_KEY=xyz
sudo -E ./run.sh
```

## Access your lambda

Point your browser at:  http://localhost:2015/

You should see the same calc form you saw in the 1st step above
when running calcweb locally on :8080

## How does this work?

* Caddy accepts the HTTP request and hands it off to the `awslambda` plugin
* The `awslambda` plugin forms a JSON payload based on the HTTP request and
invokes the Lambda specified in the `Caddyfile`
* AWS Lambda invokes the `caddy_web` function with the JSON payload
* `lambazon` Go lib creates `net/http` request/response types and calls `ServeHTTP` on the web app
* `calcweb` code runs normally, blissfully unaware it's running in Lambda
* `lambazon` creates a JSON response payload using the envelope format specified by the Caddy `awslambda` plugin
* `awslamda` plugin receives the JSON response from AWS and translates it into a HTTP response
* Caddy sends HTTP response back to browser

## Limitations

* Since Apex only deploys a single binary, no local static files will be available
  * To workaround, consider using a tool like https://github.com/jteeuwen/go-bindata
