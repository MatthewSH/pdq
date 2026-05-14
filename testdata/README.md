# Test Data
This is a directory that contains test images to validate that our implementation is as close to specifications from Meta.

## Fetching Data
```shell
go run ./internal/cmd/fetch-testdata
```

Images are downloaded from Meta's ThreatExchange repository:
https://github.com/facebook/ThreatExchange/tree/main/pdq/data

## Running integration tests

```
go test -tags integration
```