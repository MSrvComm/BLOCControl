# A simple control plane for MiCo

***This is still largely a work in progress***

The overall idea is to build a kubernetes `DaemonSet` that watches kubernetes service endpoints and also is a pubsub server (websocket/grpc) using protobuf. This can provide targeted information to endpoints and can also implement custom autoscalers.

## The controller

The current version of the controller is implemented in the `restEpWatcher` folder.

## Watch endpoints

The `endpoint_watcher.go` implements an endpoint watcher.

## The pubsub server

The `web/srv.go` and other folders are about experimenting and refining the pubsub component.
