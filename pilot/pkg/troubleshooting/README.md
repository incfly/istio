# Troubleshooting Setup Guide


## Setup

server

```shell
source setup.sh && server
```

proxy agent

```shell
# first proxy window
source ./setup.sh&& agent -i 'proxy1'
# second one
source ./setup.sh&& agent -i 'proxy2'
```

CLI window

```shell
# diff order of executions.
source ./setup.sh && cli  -s  ''
source ./setup.sh && cli  -s  'proxy2'
source ./setup.sh && cli  -s  'random'
```

## TODO

In order...

1. actual respect the requestId. incrementing.
1. maybe tracking map when connection is lost?
1. HTTP libraries for sending request to config dump interface.

### DONE

1. Handle TODO for long running request streaming.
   1. Done. Already in separate go routine.
1. More than one proxy streaming, be able to stream response, linked replies.
   1. add selector in proto.
   1. hardcoded syntax for now, "proxy1" prefix, "all" for the rest.
   1. cli code change for the stream effect.

1. Test Cases
  1. single 1,1,1 working.
  1. same request sending again.
  1. same request concurrently, should be able to served concurrently. (not currently, assuming...?)
  1. missing proxy id, not find.
  1. multiple find.
