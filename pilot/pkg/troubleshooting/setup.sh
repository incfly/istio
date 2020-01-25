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

# https://kubernetes.io/docs/tasks/administer-cluster/access-cluster-api/#without-kubectl-proxy
# one off setup isntructions.
# kubectl create sa echo-sa
# kubectl create clusterrolebinding cluster-admin-binding-echosa  --clusterrole=cluster-admin --user=system:serviceaccount:default:echo-sa
apiserver-foo() {
  export CLUSTER_NAME="gke_jianfeih-test_us-central1-a_istio-dev"
  TOKEN=$(kubectl get secrets -o jsonpath="{.items[?(@.metadata.annotations['kubernetes\.io/service-account\.name']=='echo-sa')].data.token}"|base64 --decode)
  APISERVER=$(kubectl config view -o jsonpath="{.clusters[?(@.name==\"$CLUSTER_NAME\")].cluster.server}")
  curl -X GET "$APISERVER/apis/echo.example.com/v1alpha1/foo/bar?sleep=20" --header "Authorization: Bearer $TOKEN" --insecure   -N
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