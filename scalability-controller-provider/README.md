# Scalability Controller Provider
## Overview

`scalability-controller-provider` is a [Crossplane](https://crossplane.io/) Provider 
designed to act as a scalability controller within an IoT ecosystem. 

The controller monitors and analyzes the performance metrics of existing MQTT Broker 
instances. Based on the incoming data load from producers, it automatically 
orchestrates the lifecycle of additional broker instances across a distributed 
inter-cluster environment, ensuring seamless horizontal scaling and high availability.

### Key Capabilities
* **Traffic Analysis:** Monitors MQTT broker load and producer demand.
* **Inter-cluster Orchestration:** Provisions managed resources across multiple Kubernetes clusters using Crossplane.
* **Dynamic Load Distribution:** Balances incoming IoT telemetry by scaling the messaging infrastructure horizontally.
