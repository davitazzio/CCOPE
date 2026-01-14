# CCOPE: IoT Automation Infrastructure Platform

This repository contains a modular, layer-based platform for **IoT Automation Systems**. Built on top of [Crossplane](https://crossplane.io/), the platform automates the entire infrastructure lifecycle—from bare-metal processes to distributed workload orchestration.

---

## Platform Architecture

The system is organized into functional layers. Each layer builds upon the abstractions provided by the one below it, creating a seamless flow from physical resources to high-level applications.

```text
 ┌─────────────────────────────────────────────────────────────────┐
 │                      Workload Layer                             │
 │           (Federated Learning & Pipeline Delivery)              │
 ├───────────────────────────────┬─────────────────────────────────┤
 │  [federated-pipeline-provider]│ [federatedlearning-processes]   │
 └───────────────▲───────────────┴───────────────▲─────────────────┘
                 │                               │
 ┌───────────────┴───────────────────────────────┴─────────────────┐
 │                      Services Layer                             │
 │           (Scalability, Fault Tolerance & Routing)              │
 ├───────────────────────┬───────────────────────┬─────────────────┤
 │ [scalability-provider]│    [proxy-provider]   │ [topic-provider]│
 └───────────────▲───────┴───────────────▲───────┴─────────▲───────┘
                 │                       │                 │
 ┌───────────────┴───────────────────────────────┴─────────────────┐
 │                      Resource Layer                             │
 │                (Lifecycle & Base Abstractions)                  │
 ├─────────────────────────────────────────────────────────────────┤
 │                 [baremetal-process-provider]                    │
 └─────────────────────────────────────────────────────────────────┘

