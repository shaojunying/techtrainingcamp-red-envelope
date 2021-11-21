# Consumer

## 容器化部署方法

```bash
docker build -t cr-cn-beijing.volces.com/group1/consumer1:v0.3 -f Dockerfile1 .
docker push cr-cn-beijing.volces.com/group1/consumer1:v0.3
docker build -t cr-cn-beijing.volces.com/group1/consumer2:v0.3 -f Dockerfile2 .
docker push cr-cn-beijing.volces.com/group1/consumer2:v0.3
```