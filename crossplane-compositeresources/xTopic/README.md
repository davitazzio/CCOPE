# xTopic CompositeResource 
xTopic is a composite resource that manages all the custom resources that are needed
to start a Topic for a pub-sub system. 

In this implementation, a MQTT broker is used, with the [EMQX](https://www.emqx.com/en) implementation. 
The broker starts on a cluster throw the EMQX operator. The cluster that hosts the broker must have the [EMQX operator installed](https://docs.emqx.com/en/emqx-operator/latest/getting-started/getting-started.html#emqx-open-source-5). 


## Prepare the testbed using Kind

The cluster that host the broker must expose the Kubernetes API server on an external IP. The Kind Cluster must be created applying the config file: 

```yaml
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
networking:
  # Address of the host machine of the kind cluster.
  apiServerAddress: "192.168.17.25"
  # By default the API server listens on a random open port.
  # You may choose a specific port but probably don't need to in most cases.
  # Using a random port makes it easier to spin up multiple clusters.
  apiServerPort: 6443
```

Once the cluster is created, you can obtain the kubeconfig by running: 
```bash
kind get kubeconfig | base64
```

Now the kubernetes-provider, created by the the XRD, is able to accept to the cluster APIs. 
The folder Examples, contain the file [topic-test.yaml](https://gitlab.com/MMw_Unibo/platformeng/crossplane-compositeresources/-/blob/main/xTopic/Examples/topic-test.yaml?ref_type=heads) with parameters. 

```yaml
apiVersion: slices.crossplane.io/vlalphal
kind: xtopic
metadata:
  name: topicxr #insert here the resource name 
spec:

  id: "provatopicprovider" #insert here the topic name
  version: v1alpha1
  host: "dtazzioli-edge.cloudmmwunibo.it" # insert here the broker host IP
  username: "77516558eb5b721f" # insert here the topic broker API username
  password: "o9CYsHQ8VPTQ26ZkiylQ8c6Obd61nNEetQi1DKEDBJzC" # Insert here the API key
  degradedtreshold: "20"  # the number of messages in the queue that define the topic as degraded
  remotekubeconfig:   # insert here the kubeconfig command output
```
