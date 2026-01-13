# Crossplane CompositeResources

## create cluster kind
- exposing the cluster at 192.168.17.25:6443.

- forwarding port 30001 to every address for services.

```yaml
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
networking:
  # Address of the host machine of the kind cluster.
  apiServerAddress: "dtazzioli-cloud.cloudmmwunibo.it"
  # By default the API server listens on a random open port.
  # You may choose a specific port but probably don't need to in most cases.
  # Using a random port makes it easier to spin up multiple clusters.
  apiServerPort: 6443
nodes:
  - role: control-plane
    extraPortMappings:
      - containerPort: 30001
        hostPort: 30001
        listenAddress: "0.0.0.0" # Optional, defaults to "0.0.0.0"
        protocol: tcp # Optional, defaults to tcp

```

Once the cluster is created, you can obtain the kubeconfig by running: 
```bash
kind get kubeconfig | base64
```

Now the kubernetes-provider, created by the the XRD, is able to accept to the cluster APIs. 

# install crossplane
```shell
  helm install crossplane \
--namespace crossplane-system \
--create-namespace crossplane-stable/crossplane 

```
# exposing cluster:
expose cluster at port 8080.

```shell
  kubectl proxy --port=8080 --address "0.0.0.0" --accept-hosts ".*"
```
Access it with rest at the selected port

## provider to use kubernetes resources in crossplane composite resources

```shell
crossplane xpkg install provider xpkg.upbound.io/crossplane-contrib/provider-kubernetes:v0.13.0
```
```shell
SA=$(kubectl -n crossplane-system get sa -o name | grep provider-kubernetes | sed -e 's|serviceaccount\/|crossplane-system:|g')
kubectl create clusterrolebinding provider-kubernetes-admin-binding --clusterrole cluster-admin --serviceaccount="${SA}"
kubectl apply -f config-in-cluster.yaml
```

# config-in cluster to control local cluster
```yaml
apiVersion: kubernetes.crossplane.io/v1alpha1
kind: ProviderConfig
metadata:
  name: default
spec:
  credentials:
    source: InjectedIdentity
```
# config-in cluster to control local cluster
```yaml
apiVersion: v1
data: 
    kubeconfig: ${kubeconfig}
kind: Secret
metadata:
    name: kubernetes-provider-secrets
    namespace: crossplane-system
type: Opaque
---
apiVersion: kubernetes.crossplane.io/v1alpha1
kind: ProviderConfig
metadata:
    annotations:
    name: kubernetes-provider-config
spec:
    credentials:
    secretRef:
        key: kubeconfig
        name: kubernetes-provider-secrets
        namespace: crossplane-system
    source: Secret
```
