# scheduler extender

## 拷贝 extender 配置文件

```bash
# vim config/scheduler-extender-config.yaml
...
extenders:
  - urlPrefix: "http://x.x.x.x:8000" # 修改这个ip为你的 extender server ip
...

# cp config/scheduler-extender-config.yaml /etc/kubernetes/scheduler-extender-config.yaml
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
      - mountPath: /etc/kubernetes/scheduler-extender-config.yaml # 添加这一行
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

# kubectl apply -f deploy/deployment.yaml

# kubectl apply -f deploy/service.yaml
```

```bash
# kubectl get pod -n default | grep scheduler-extender
scheduler-extender-58c99bf48f-25kd6   1/1     Running   0          79m

# kubectl get svc -n default | grep scheduler-extender
scheduler-extender   NodePort    10.110.24.202   <none>        8000:31234/TCP   75m

# kubectl logs -f scheduler-extender-58c99bf48f-25kd6 -n default

```

## 手动 启动 extender server

```bash
# go mod tidy

# go run main.go
```