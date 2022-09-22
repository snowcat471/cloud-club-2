# Cloud Club 2기 프로젝트

## 2주차 목표

- k8s resource(Pod, Replicaset, Deployment, Service 등)를 만들어보고 컨트롤 해보기

## 1. minikube 세팅 (MacOS)

~~~bash
# minikube 및 kubectl 설치
$ brew install minikube kubectl

# minikube 시작
$ minikube start
~~~

## 2. Pod, Service 배포

- 1주차에 만든 도커이미지를 수정하여 v1, v2 2가지로 실습

### Imperative 

- Pod 생성

~~~bash
# Pod 생성
$ kubectl run --image=snowcat471/simple-app:v1 -l app=web-app web-app
pod/web-app created

# Pod 확인
$ kubectl get pod
NAME     READY   STATUS    RESTARTS   AGE
web-app  1/1     Running   0          24s
~~~

- Service(ClusterIP) 생성

~~~bash
# Service 생성
$ kubectl expose pod web-app --port=80 --target-port=3000
service/web-app exposed

# Service 확인
$ kubectl get svc
NAME         TYPE        CLUSTER-IP      EXTERNAL-IP   PORT(S)   AGE
kubernetes   ClusterIP   10.96.0.1       <none>        443/TCP   51m
web-app      ClusterIP   10.111.201.46   <none>        80/TCP    9s

# endpoint 연결 확인
$ describe svc web-app
Name:              web-app
Namespace:         default
Labels:            app=web-app
Annotations:       <none>
Selector:          app=web-app
Type:              ClusterIP
IP Family Policy:  SingleStack
IP Families:       IPv4
IP:                10.111.201.46
IPs:               10.111.201.46
Port:              <unset>  80/TCP
TargetPort:        3000/TCP
Endpoints:         172.17.0.3:3000 # <pod ip>:<port>
Session Affinity:  None
Events:            <none>
~~~

- Service(NodePort) 생성

~~~bash
# 생성
$ kubectl expose pod web-app --type=NodePort --port=8000 --target-port=3000
service/web-app exposed

# 확인
$ kubectl get svc
NAME         TYPE        CLUSTER-IP       EXTERNAL-IP   PORT(S)          AGE
kubernetes   ClusterIP   10.96.0.1        <none>        443/TCP          62m
web-app      NodePort    10.108.234.178   <none>        8000:31150/TCP   29s

$ kubectl describe svc web-app
Name:                     web-app
Namespace:                default
Labels:                   app=web-app
Annotations:              <none>
Selector:                 app=web-app
Type:                     NodePort
IP Family Policy:         SingleStack
IP Families:              IPv4
IP:                       10.108.234.178
IPs:                      10.108.234.178
Port:                     <unset>  8000/TCP
TargetPort:               3000/TCP
NodePort:                 <unset>  31150/TCP
Endpoints:                172.17.0.3:3000
Session Affinity:         None
External Traffic Policy:  Cluster # <pod ip>:<port>
Events:                   <none>
~~~

~~~
# minikube ip 확인 
$ minikube ip
192.168.49.2

# Service access 확인
$ curl 192.168.49.2:31150 # 접근이 안되는 것을 확인

# "minikube service --url web-app" 명령어를 통해 서비스 접근 주소 생성해야 함
http://127.0.0.1:62336

# Service access 재확인
$ curl 127.0.0.1:62336
Hello v1 # 응답 확인
~~~

- Service(LoadBalancer) 

~~~bash
# 생성
$ kubectl expose pod web-app --type=LoadBalancer --port=80 --target-port=3000
service/web-app exposed

$ kubectl get svc
NAME         TYPE           CLUSTER-IP       EXTERNAL-IP   PORT(S)        AGE
kubernetes   ClusterIP      10.96.0.1        <none>        443/TCP        71m
web-app      LoadBalancer   10.102.170.155   <pending>     80:30743/TCP   18s

### minikube에서 실습중이므로 EXTERNAL-IP가 계속 <pending> 상태로 남아있음
### AWS, GCP, Azure 등 클라우드 서비스 사용시 외부아이피를 부여받아 외부에서 접근 가능
~~~

### Declarative

yaml 파일을 이용하여 선언적으로 관리<br>
[Kubernetes Documentation](https://kubernetes.io/docs/home/)을 참고하여 직접 yaml 파일 생성하거나, -o yaml 옵션을 통해 yaml 파일을 만들 수 있음<br>

- Pod

~~~bash
# yaml 파일 생성
$ kubectl run --image=snowcat471/simple-app:v1 -l app=web-app web-app --dry-run=client -o yaml > web-app.yaml
### 생성한 yaml 파일에서 필요없는 필드는 지우고 정리
~~~

~~~yaml
apiVersion: v1
kind: Pod
metadata:
  labels:
    app: web-app
  name: web-app
spec:
  containers:
  - image: snowcat471/simple-app:v1
    name: web-app
~~~

~~~bash
# 파일을 통해 pod 생성
$ kubectl apply -f web-app.yaml
pod/web-app created

# 확인
$ kubectl get pod
NAME      READY   STATUS    RESTARTS   AGE
web-app   1/1     Running   0          67s
~~~

- Service

~~~bash
# ClusterIP
$ kubectl expose pod web-app --port=80 --target-port=3000 --dry-run=client -o yaml > web-app-clusterip.yaml
$ kubectl create svc clusterip web-app --tcp=80:3000 --dry-run=client -o yaml > web-app-clusterip.yaml

# NodePort
$ kubectl expose pod web-app --type=NodePort --port=8000 --target-port=3000 --dry-run=client -o yaml > web-app-nodeport.yaml
$ kubectl create svc nodeport web-app --tcp=8000:3000 --dry-run=client -o yaml > web-app-nodeport.yaml

# LoadBalancer
$ kubectl expose pod web-app --type=LoadBalancer --port=80 --target-port=3000 --dry-run=client -o yaml > web-app-loadbalancer.yaml
$ kubectl create svc loadbalancer web-app --tcp 80:3000 --dry-run=client -o yaml > web-app-loadbalancer.yaml
~~~

~~~yaml
# web-app-clusterip.yaml
apiVersion: v1
kind: Service
metadata:
  labels:
    app: web-app
  name: web-app
spec:
  ports:
  - port: 80
    protocol: TCP
    targetPort: 3000
  selector:
    app: web-app
~~~

~~~yaml
# web-app-nodeport.yaml
apiVersion: v1
kind: Service
metadata:
  labels:
    app: web-app
  name: web-app
spec:
  ports:
  - port: 8000
    protocol: TCP
    targetPort: 3000
  selector:
    app: web-app
  type: NodePort
~~~

~~~yaml
# web-app-loadbalancer.yaml
apiVersion: v1
kind: Service
metadata:
  labels:
    app: web-app
  name: web-app
spec:
  ports:
  - port: 80
    protocol: TCP
    targetPort: 3000
  selector:
    app: web-app
  type: LoadBalancer
~~~

~~~bash
# 생성
$ kubectl apply -f <file path>
~~~

### Pod 형태로 배포시 문제점

- 새로운 버전 배포시 직접 새버전의 Pod를 추가하고, 구버전의 Pod를 삭제해야함
- Scale Out 불가능

## Replicaset

~~~yaml
apiVersion: apps/v1
kind: ReplicaSet
metadata:
  name: web-app
  labels:
    app: web-app
spec:
  replicas: 3
  selector:
    matchLabels:
      app: web-app
  template:
    metadata:
      labels:
        app: web-app
    spec:
      containers:
      - name: web-app
        image: snowcat471/simple-app:v1
~~~

~~~bash
# 생성
$ kubectl apply -f rs.yaml 
replicaset.apps/web-app created

# 확인
$ kubectl get rs
NAME      DESIRED   CURRENT   READY   AGE
web-app   3         3         3       22s

$ kubectl get pod
NAME            READY   STATUS    RESTARTS   AGE
web-app-9ztw6   1/1     Running   0          43s
web-app-k7tpf   1/1     Running   0          43s
web-app-kdtmq   1/1     Running   0          43s

# scale
$ kubectl scale rs web-app --replicas=5
replicaset.apps/web-app scaled

$ kubectl get pod
NAME            READY   STATUS    RESTARTS   AGE
web-app-8r4jm   1/1     Running   0          19s
web-app-9ztw6   1/1     Running   0          102s
web-app-k7tpf   1/1     Running   0          102s
web-app-kdtmq   1/1     Running   0          102s
web-app-mv4vr   1/1     Running   0          19s

# pod 삭제해보기
$ kubectl delete pod web-app-kdtmq 
pod "web-app-kdtmq" deleted

$ kubectl get pod
NAME            READY   STATUS    RESTARTS   AGE
web-app-8r4jm   1/1     Running   0          2m19s
web-app-9ztw6   1/1     Running   0          3m43s
web-app-b8b62   1/1     Running   0          8s
web-app-k7tpf   1/1     Running   0          3m43s
web-app-mv4vr   1/1     Running   0          2m19s
~~~

### Replica Set 배포의 문제점

- Pod와 마찬가지로 새로운 버전으로의 업데이트가 지원되지 않음

## Deployment

~~~yaml

~~~
