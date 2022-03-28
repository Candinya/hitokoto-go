# Hitokoto-Go

> Hitokoto 的 Golang 版本

## 关于

- Q: 为什么要写这个版本？

  A: 没什么，就是玩儿 🥰

- Q: 这个版本希望达到的目标？

  A: 更强大的性能，更强大的功能和扩展性

- Q: 为什么日志是英文？

  A: 非 ASCII 字符可能在控制台输出时会出现不同的表现，为避免问题，全部使用英文日志。

- Q: 网易云音乐呢？

  A: [我不想吃 DMCA](https://github.com/github/dmca/blob/5e1f9b145d35dac15a012a725cf1a399c0a17f16/2021/06/2021-06-21-netease.md)

## 使用

### 写在前面

1. 我们推荐使用 Docker 环境来运行本容器，使用 docker-compose 进行管理；您也可以选择编译源码，或是下载二进制文件运行的方案。
2. 您需要使用环境变量来指定数据库（我们使用的是 postgresql ）和 redis 的连接字符串， docker 环境下的已经默认配置完成，二进制环境下的您需要手动设置。

### 导入数据

在使用之前，您需要先初始化并导入数据库。

我们使用了 [一言社区的官方句子库](https://github.com/hitokoto-osc/sentences-bundle) ，当前支持的协议版本是 v1.0.0 。请注意，如果协议版本有更新，请提示我们跟进新的支持。

您需要下载该句子仓库至与本 ReadMe.md 同级的目录（如果是 docker-compose 方案，则是与 docker-compose 文件同级目录），以便程序运行时读取数据。如果目录不正确，可能会找不到文件，导致导入失败。

#### docker-compose 方式

1. 取消注释 `docker-compose.yml` 文件中的第 20、21 行关于卷映射的描述。
2. 运行 `docker-compose run --rm app hitokoto-go --import` 命令，初始化 docker 运行环境，并导入数据。如果系统提示无法连接到数据库，可能是数据库还未启动，您可稍后重试。
3. 由于在正式运行环境中不需要导入文件的映射，所以推荐您在数据导入工作完成后重新恢复第一步中取消掉的注释。

#### 二进制文件方式

使用二进制文件方案的用户可以在启动二进制文件时附带 `--import` 参数来启动导入数据模式。

数据导入完成后，系统会提示 `All data imported! enjoy :D` ，此时您可不带参数地启动程序，以提供服务。

## 开发

您可使用 `docker-compose.dev.yml` 文件中指定的参数来启动测试环境所需要的依赖。

当然，您需要手动设置以下环境变量：

```
POSTGRES_CONNECTION_STRING=postgres://hitokoto:hitokoto@localhost:5432/hitokoto
REDIS_CONNECTION_STRING=redis://localhost:6379/0
```

## 开发进度

### 测试脚本

- [ ] 逻辑验证
- [ ] 请求模拟

### 导入导出

- [x] 数据导入与更新
- [ ] 数据导出（JSON格式）

### 读取请求

- [x] 一言分类与请求条件
- [x] 句子长度与请求条件
- [ ] 跨表随机/权重随机
- [ ] 基于请求数量/喜欢数量的一言选取方式
- [ ] 全文搜索

### 写入请求

- [ ] 提交
- [ ] 请求统计，喜欢统计等
- [ ] 管理员操作

### 其他

其他的好想法？欢迎随时开一个 Issue 来帮助我们变得更棒！

## 性能

当前版本的性能极其拉跨，由于使用了数据库随机算法，所以花费了非常高昂的计算费用和访问开销。将着重针对随机算法部分进行优化。

## 开发团队

- Nya Candy

## 非常感谢

- [一言开源社区](https://github.com/hitokoto-osc)
