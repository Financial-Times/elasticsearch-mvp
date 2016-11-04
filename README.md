# POC for using ElasticSearch 

Small application, which checks the health of an AES cluster.

:warning: The AWS SDK for Go [does not currently include support for ES data plane api](https://github.com/aws/aws-sdk-go/issues/710), but the Signer is exposed since v1.2.0.

The taken approach to access ES:
- Use the v4.Signer provided by the amazon-go-sdk
- Create an HTTP client wrapping all the request with Amazon signer (github.com/sha1sum/aws_signing_client)
- Use https://github.com/olivere/elastic library to any ES request, after passing on the above created client

## How to run

```
go get -u github.com/Financial-Times/elasticsearch-mvp
go build
./elasticsearch-mvp --aws-access-key={access key} --aws-secret-access-key={secret key}
```
Please provide the elasticsearch endpoint, region and the host you want to run the app on.

## Available endpoints:

### localhost:8080/__health

Provides the standard FT output indicating the clusters health.

### localhost:8080/__health-details

Provides a detailed health status from ES cluster. It matches the response from [elasticsearch-endpoint/_cluster/health](https://www.elastic.co/guide/en/elasticsearch/reference/current/cluster-health.html)