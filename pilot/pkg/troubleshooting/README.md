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

### Multi cli req setup

multiple cli debugging request, request id is respected correctly.

```shell
source ./setup.sh && server | tee server.log
source ./setup.sh&& agent -i 'proxy1' | tee client.log
source ./setup.sh && multiclient
```

### Kubernetes Workflow

Deploy service

```shell
kubectl apply -f ./troubleshooting.yaml
```

Observing, each in separate terminal window.

```shell
kpfn istio-system $(kpidn istio-system -lapp=ts-server) 8000
klo -lapp=httpbin -c istio-proxy -f
klon istio-system -lapp=ts-server
source ./setup.sh && cli  -s  'sidecar' -r $(date)
```

Rebuild

```shell
source setup.sh; docker-build
source setup.sh; deploy
```

## Aggregation Layer

1. long duration req is possible
1. streaming is possible (curl -N option)
1. using customized path to be k8s-noic

```shell
source setup.sh && apiserver-foo
 ```


## Notes

2020/1/2 discussion

- api server only do uri based authorization, not payload, which can touch multi namespace.
- authorization check can be done at single entry replica
- need a design for request not always fan out
  - diff metadata
  - diff grpc
  - different port.

## TODO

In right order.

1. Propogate sidecar actual id to the proxy id.
1. Tie response to the proxy id.
1. GC management for the proxy id in the map, when connection is lost.
1. request log scope formatting polishing, requets id, proxy id, as base context.

### DONE

1. GC requestMap, proxyInfo map.
  1. without this, unable to invoke twice. need manual specifying request id. wrong

1. Things working E2E.
  1. Standalone separate server.
  1. Client compile into istiod agent.
  1. Dial against server. Hardcoded.
  1. Port forwarding tryout e2e.

1. actual respect the requestId. incrementing.

why need the request id? because agent sharing the same stream with different istoctl session.
  1. cli slow req1, 5 secs, hold.
  1. cli req2, 5 secs, might get the req1 s response?

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


## QA & Ideas

Ideas

- Make istioctl anlyzing functionatlity of envoy config can consume streaming stdin:

  ```shell
  istioctl troubleshooting --selector='app=foo' | istioctl proxy-config
  ```

- Does gRPC context auto populate metadata for the tracing? Or code required?