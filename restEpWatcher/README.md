# EPWatcher

This controller watches for endpoints defined for each service in the cluster. It is also listening on port `62000` of the pod. A request of the format `http://epwatcher:62000/testapp-svc-2/` will return endpoints for the `testapp-svc-2` service.

## Building the watcher

```bash
sudo docker build -t ratnadeepb/epwatcher:latest .
sudo docker push ratnadeepb/epwatcher:latest
```
