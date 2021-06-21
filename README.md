<h1 align="center">Fusion Ops</h1>

<div align="center">
 Fusion 运营支撑系统
<br/>

</div>

# 基于Makefile运行
$ make start

# OR 使用go命令运行
$ go run ./main.go web -c ./config_file/config.toml -m ./config_file/model.conf --menu ./config_file/menu.yaml
```

> 启动成功之后，可在浏览器中输入地址进行访问：[http://127.0.0.1:10088/swagger/index.html](http://127.0.0.1:10088/swagger/index.html)

## 生成`swagger`文档

```bash
# 基于Makefile
make swagger

# OR 使用swag命令
swag init --parseDependency --generalInfo ./main.go --output ./swagger
```

## 重新生成依赖注入文件

```bash
# 基于Makefile
make wire

# OR 使用wire命令
wire gen ./app/app
```
