# Troubleshooting Setup Guide


## Setup

First window

```shell
source setup.sh && server
```

Second window

```shell
source setup.sh && agent
```

Third window

```shell
source setup.sh && cli
```

## TODO

In order...

1. Handle TODO for long running request streaming.
1. More than one proxy streaming, be able to stream response, linked replies.
   1. use reserved word "all" for now.
1. maybe tracking map when connection is lost?
1. HTTP libraries for sending request to config dump interface.
1. Test.