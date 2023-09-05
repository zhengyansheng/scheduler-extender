# scheduler extender

## 拷贝 extender 配置文件

```bash
# cp config/scheduler-extender-config.yaml /etc/kubernetes/
```

## 修改 manifests/scheduler 的 yaml文件

```bash
# vim /etc/kubernetes/manifests/kube-scheduler.yaml
```

```yaml
apiVersion: v1
kind: Pod
...
spec:
  containers:
    - command:
        - kube-scheduler
        - --authentication-kubeconfig=/etc/kubernetes/scheduler.conf
        - --authorization-kubeconfig=/etc/kubernetes/scheduler.conf
        - --bind-address=127.0.0.1
        - --kubeconfig=/etc/kubernetes/scheduler.conf
        - --leader-elect=true
        - --config=/etc/kubernetes/scheduler-extender-config.yaml # 添加这一行

      ......

      volumeMounts:
      - mountPath: /etc/kubernetes/scheduler.conf
        name: kubeconfig
        readOnly: true
      - mountPath: scheduler-extender-config.yaml # 添加这一行
        name: scheduler-extender
        readOnly: true

      ......
  volumes:
    - hostPath:
        path: /etc/kubernetes/scheduler.conf
        type: FileOrCreate
      name: kubeconfig
    - hostPath: # 添加这一行
        path: /etc/kubernetes/scheduler-extender-config.yaml
        type: FileOrCreate
      name: scheduler-extender
    ...

```


## 部署 extender server

```bash

# kubectl apply -f manifests/deployment.yaml

# kubectl apply -f manifests/service.yaml
```

```text

## 启动 extender server

```bash
# go mod tidy

# go run main.go
```

```text
[root@master multi-scheduler-webhook]# go run main.go
[GIN-debug] [WARNING] Creating an Engine instance with the Logger and Recovery middleware already attached.

[GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
 - using env:	export GIN_MODE=release
 - using code:	gin.SetMode(gin.ReleaseMode)

[GIN-debug] GET    /ping                     --> main.main.func1 (3 handlers)
[GIN-debug] POST   /scheduler/extender/filter --> main.main.func2 (3 handlers)
[GIN-debug] POST   /scheduler/extender/prioritize --> main.main.func3 (3 handlers)
[GIN-debug] [WARNING] You trusted all proxies, this is NOT safe. We recommend you to set a value.
Please check https://pkg.go.dev/github.com/gin-gonic/gin#readme-don-t-trust-all-proxies for details.
[GIN-debug] Listening and serving HTTP on 0.0.0.0:8000
[GIN] 2023/08/30 - 14:57:57 | 404 |         952ns |     10.112.0.20 | POST     "/xx/filter"
I0830 14:58:29.327981    3458 main.go:37] filter localstorage pods on nodes: [master master2 node1]
&v1.ExtenderFilterResult{
    Nodes:                      (*v1.NodeList)(nil),
    NodeNames:                  &[]string{"master", "master2", "node1"},
    FailedNodes:                {},
    FailedAndUnresolvableNodes: {},
    Error:                      "",
}[GIN] 2023/08/30 - 14:58:29 | 200 |    3.706826ms |     10.112.0.20 | POST     "/scheduler/extender/filter"
I0830 14:58:29.329429    3458 main.go:60] scoring nodes [master master2 node1]
I0830 14:58:29.329447    3458 main.go:71] score localstorage pods on nodes: [{master 1} {master2 7} {node1 7}]
&v1.HostPriorityList{
    {Host:"master", Score:1},
    {Host:"master2", Score:7},
    {Host:"node1", Score:7},
}[GIN] 2023/08/30 - 14:58:29 | 200 |      273.94µs |     10.112.0.20 | POST     "/scheduler/extender/prioritize"
```