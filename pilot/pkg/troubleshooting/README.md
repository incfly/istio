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

1. More than one proxy streaming, be able to stream response, linked replies.
   1. add selector in proto
   1. hardcoded syntax for now, "proxy1" prefix, "all" for the rest.
   1. cli code change for the stream effect.
1. actual respect the requestId. incrementing.
1. maybe tracking map when connection is lost?
1. HTTP libraries for sending request to config dump interface.


1. Handle TODO for long running request streaming.
   1. Done. Already in separate go routine.


1. Test Cases
  1. single 1,1,1 working.
  1. same request sending again.
  1. same request concurrently, should be able to served concurrently. (not currently, assuming...?)
  1. missing proxy id, not find.
  1. multiple find.
