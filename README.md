# istio-microservice-example

This repo provides a `kubectl apply -f kubernetes/` example to show two ping/pong
services within istio.

This is aimed at someone who has istio installed but wants to generate real traffic.

![flow](res/flow.gif)

![1](res/1.png)

## Overview of code

A single file that does both server/client ping/pong via `protocolbuffers/message` type

The magic that makes istio work is `kubernetes/*/istio-virtualservice`

![2](res/2.png)