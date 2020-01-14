#!/bin/bash

proto() {
  protoc --go_out=plugins=grpc:. api/service.proto
}

# build() {
#   go build -o ts-server 
# }

docker-build() {
  pushd cmd/server
  go build -o ts-server main.go
  docker build . -t gcr.io/jianfeih-test/ts-server:0108a
  docker push gcr.io/jianfeih-test/ts-server:0108a
  popd

  pushd "${GOPATH}/src/istio.io/istio"
  export TAG="0108a" HUB="gcr.io/jianfeih-test"
  make docker.proxyv2 && docker push "gcr.io/jianfeih-test/proxyv2:${TAG}"
  popd
}

deploy() {
  krmpo -nistio-system -lapp=ts-server
  krmpo -lapp=httpbin
  k apply -f ./troubleshooting.yaml
  k logs -lapp=httpbin -c istio-proxy -f
}

server() {
  go run ./cmd/istiod "$@"
}

agent() {
  go run ./cmd/agent "$@"
}

cli() {
  go run ./cmd/cli "$@"
}

genhttp2cert() {
  openssl req -newkey rsa:2048 -nodes -keyout server.key -x509 -days 365 -out server.crt
}


# out of order execution
# starting with 1
# 2019-12-30T20:21:39.999175Z	info	respose is {response-proxy1-cli-req-5 {} [] 0}
# starting with 10
# 2019-12-30T20:21:40.425842Z	info	respose is {response-proxy1-cli-req-3 {} [] 0}
# starting with 2
# 2019-12-30T20:21:42.044724Z	info	respose is {response-proxy1-cli-req-6 {} [] 0}
# starting with 3
# 2019-12-30T20:21:39.142941Z	info	respose is {response-proxy1-cli-req-8 {} [] 0}
# starting with 4
# 2019-12-30T20:21:40.268722Z	info	respose is {response-proxy1-cli-req-9 {} [] 0}
# starting with 5
# 2019-12-30T20:21:39.286300Z	info	respose is {response-proxy1-cli-req-10 {} [] 0}
# starting with 6
# 2019-12-30T20:21:40.392495Z	info	respose is {response-proxy1-cli-req-2 {} [] 0}
# starting with 7
# 2019-12-30T20:21:42.969102Z	info	respose is {response-proxy1-cli-req-4 {} [] 0}
# starting with 8
# 2019-12-30T20:21:39.067125Z	info	respose is {response-proxy1-cli-req-7 {} [] 0}
# starting with 9
# 2019-12-30T20:21:39.366402Z	info	respose is {response-proxy1-cli-req-1 {} [] 0}
multiclient() {
  rm -rf output/ && mkdir output/
  for i in `seq 1 5`; do
    # Bash parallel in for loop is noter deterministically... will be out of order by itself...
    # fixing by adding param explicitly as part of rpc.
    # TODO: understand bash behavior...
    echo "req $i started"; echo "starting with iii-$i" > output/$i.txt;
    source ./setup.sh && cli  -s  'p' -r "iii-$i" 2>&1 | grep 'respose is' >> output/$i.txt & ;
  done
  # wait for some time all finishes.
  sleep 10
  cat output/* > somefile.txt
}