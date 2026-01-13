# MQTT Test Suite: Dummy Producer & Consumer

This directory contains a specialized Go-based toolset designed to stress-test and validate the scalability and routing logic of the `scalability-controller` and `proxy-provider`.

It consists of two main components:
1.  **Dummy Producer**: Generates synthetic IoT telemetry at a configurable rate.
2.  **Dummy Consumer**: Receives and processes messages, simulating real-world consumption latency.

## Architecture Overview

These tools are designed to simulate varying workloads to trigger scaling events in the inter-cluster infrastructure. By adjusting the message frequency, you can simulate peak loads or network congestion.



---

## Getting Started

### Prerequisites
* Go 1.21+
* Access to an MQTT Broker (e.g., EMQX, Mosquitto)

### Installation
```bash
go mod download
