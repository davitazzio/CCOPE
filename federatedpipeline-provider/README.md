# federatedpipeline provider

This provider can access kubernetes API to get data from topics and federated learning clients in order to scale the clients when the topics are cluttered.

It expects the topic provider to expose a degraded flag signaling that the topic queue is getting crammed.

## provider inputs
```yaml
ClusterAddress: string  "http://dtazzioli-edge.cloudmmwunibo.it:8080"
TopicName:		  string  "mqtttopic-1"
```
### kubernetes cluster configuration to expose API and federated clients ports

1. create the cluster forwarding the port for accessing the flclient
```shell
  kind create cluster --config=kind_config.yaml
```
2. accessing external kubernetes cluster
Run on the external cluster:
```shell
  kubectl proxy --port=8080 --address "0.0.0.0" --accept-hosts ".*"
```
Access it with rest at the selected port






