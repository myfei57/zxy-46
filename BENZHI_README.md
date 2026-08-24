基于 Go 实现的 BldgHVAC 智能楼宇暖通空调群控平台项目，一款 Web 控制服务，完成冷站机组分级加载、分区排程与 VAV 末端联动、气象与需量补偿、能耗趋势与告警管理。

## 构建与运行

构建镜像：

```bash
./build_benzhi_docker.sh
```

本地运行：

```bash
go build -mod=vendor -o bldghvac.exe ./cmd/bldghvac
./bldghvac.exe -addr :8080 -data ./data
```

打开 http://localhost:8080 可访问控制台首页，导航进入分区、冷站、VAV 与告警四个页面；JSON 接口位于 /api 前缀。

数据以文件形式持久化在 -data 指定目录下，不依赖数据库。
