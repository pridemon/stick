# Stick
Small service for saving json documents from HTTP POST into Elasticsearch.

## How to run
Example of available env vars with their default values:
```bash
docker run
    -d --restart=always \
    --name stick \
    -e STICK_ELASTIC_URL=http://localhost:9200 \  ## elasticsearch url
    -e STICK_ELASTIC_HEALTHCHECK_INTERVAL=10 \    ## check elasticsearch every 10 seconds
    -e STICK_ELASTIC_MAX_RETRIES=5 \              ## maximal amount of retries
    -e STICK_COMMIT_AMOUNT=500 \       ## commit elastic bulk query after 500 documents
    -e STICK_COMMIT_INTERVAL=10 \      ## commit elastic bulk every 10 seconds
    -e STICK_INTERNAL_BUFFER=1000 \    ## size of internal message buffer
    -e STICK_WORKERS=5 \               ## amount of parallel goroutines for commiting into elasticsearch
    -p 8080:8080 \
    ontrif/stick
```

Then send any valid json to `8080` port:
```bash
curl -XPOST localhost:8080/ -d '{"test":"me"}'
```
