# events2prom 🧹

events2prom collects events and aggregates them in memory to
publish metrics in the Prometheus exposition format.
It allows enabling and disabling aggregations in runtime
and can support high dimensional cases.

## Usage

```
$ docker run -it -p 6677:6677 -p 6678:6678/udp \
    -v "${PWD}/examples/config/events2prom.yaml":/events2prom.yaml \
    public.ecr.aws/q1p8v8z2/events2prom \
    -config /events2prom.yaml
2021/11/08 13:27:17 Listening to admin server at ":6677"...
2021/11/08 13:27:17 Enabled collection: "request_count_by_pod"
2021/11/08 13:27:17 Enabled collection: "request_latency_ms"
2021/11/08 13:27:17 Listening events at [::]:6678, let's 🧹...
```

Then, run the example to generate events:

```
$ docker run -it --network host public.ecr.aws/w6n7a8r1/events2prom-example
```

Open http://127.0.0.1:6677/metrics to see the aggregated metrics.

```
# HELP request_count_by_pod Request count by pod
# TYPE request_count_by_pod counter
request_count_by_pod{pod="pod-1e0"} 1337
request_count_by_pod{pod="pod-1ff"} 1278
request_count_by_pod{pod="pod-def"} 1280
# HELP request_latency_ms Request latency in ms distribution
# TYPE request_latency_ms histogram
request_latency_ms_bucket{pod="pod-1e0",le="0"} 3
request_latency_ms_bucket{pod="pod-1e0",le="100"} 333
request_latency_ms_bucket{pod="pod-1e0",le="200"} 683
request_latency_ms_bucket{pod="pod-1e0",le="300"} 991
request_latency_ms_bucket{pod="pod-1e0",le="+Inf"} 991
request_latency_ms_sum{pod="pod-1e0"} 268975
request_latency_ms_count{pod="pod-1e0"} 991
request_latency_ms_bucket{pod="pod-1ff",le="0"} 5
request_latency_ms_bucket{pod="pod-1ff",le="100"} 328
request_latency_ms_bucket{pod="pod-1ff",le="200"} 661
request_latency_ms_bucket{pod="pod-1ff",le="300"} 978
request_latency_ms_bucket{pod="pod-1ff",le="+Inf"} 978
request_latency_ms_sum{pod="pod-1ff"} 250721
request_latency_ms_count{pod="pod-1ff"} 978
request_latency_ms_bucket{pod="pod-def",le="0"} 1
request_latency_ms_bucket{pod="pod-def",le="100"} 324
request_latency_ms_bucket{pod="pod-def",le="200"} 640
request_latency_ms_bucket{pod="pod-def",le="300"} 984
request_latency_ms_bucket{pod="pod-def",le="+Inf"} 984
request_latency_ms_sum{pod="pod-def"} 254033
request_latency_ms_count{pod="pod-def"} 984
```

## Usage on Kubernetes

Run the following to install events2prom on a Kubernetes cluster. See the config file
to edit the config:

```
$ kubectl apply -f ./examples/config/k8s.yaml
```

Run an example program to produce events2prom events:

```
$ kubectl apply -f ./examples/config/example-app-k8s.yaml
```

Run the following command to ensure both events2prom and example application is
running:

```
$ kubectl get pods -n events2prom
NAME                 READY   STATUS    RESTARTS   AGE
events2prom-7gmr8           1/1     Running   0          5m31s
events2prom-97vw6           1/1     Running   0          5m3s
events2prom-dp9z4           1/1     Running   0          5m22s
events2prom-dvwbn           1/1     Running   0          5m17s
events2prom-example-47zq6   1/1     Running   0          32s
events2prom-example-6dr8c   1/1     Running   0          32s
events2prom-example-k4b76   1/1     Running   0          32s
events2prom-example-ltkqx   1/1     Running   0          32s
events2prom-example-s5lgq   1/1     Running   0          32s
events2prom-sxtpc           1/1     Running   0          4m58s
```

events2prom will run as a DaemonSet and will publish Prometheus metrics.
Run Prometheus to scrape the events2prom output.
