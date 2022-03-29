# Hitokoto-Go

> Go implementation of [hitokoto API](https://github.com/hitokoto-osc/hitokoto-api)

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

- [ ] 无数据库模式（直接读取本地/远程 JSON 文件）
- [x] 数据导入与更新
- [ ] 数据导出（JSON格式）

### 读取请求

- [x] 纯随机
- [x] 指定分类随机
- [x] 句子长度限制
- [x] 跨表随机/权重随机
- [ ] 基于请求数量/喜欢数量权重随机

- [ ] 类型 + ID 定位
- [ ] UUID 定位

- [ ] 搜索（用户/内容）

### 写入请求

- [ ] 提交
- [ ] 统计请求、喜欢数量等
- [ ] 管理员操作

### 其他

其他的好想法？欢迎随时开一个 Issue 来帮助我们变得更棒！

## 性能

以下数据为使用单点 Docker-Compose 部署得出，仅供参考，具体情况请以实际体验为准。

### 服务器

```shell
root@core:~/hitokoto-go# wrk -t8 -c1000 -d10s --latency http://127.0.0.1:8080
Running 10s test @ http://127.0.0.1:8080
  8 threads and 1000 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency   291.75ms  264.35ms   2.00s    82.82%
    Req/Sec   505.38    346.25     1.54k    57.34%
  Latency Distribution
     50%  184.58ms
     75%  372.49ms
     90%  675.90ms
     99%    1.22s 
  39882 requests in 10.08s, 30.00MB read
Requests/sec:   3956.95
Transfer/sec:      2.98MB
root@core:~/hitokoto-go# neofetch
       _,met$$$$$gg.          root@core
    ,g$$$$$$$$$$$$$$$P.       --------------------- 
  ,g$$P"     """Y$$.".        OS: Debian GNU/Linux 11 (bullseye) x86_64 
 ,$$P'              `$$$.     Host: KVM/QEMU (Standard PC (i440FX + PIIX, 1996) pc-i440fx-4.0) 
',$$P       ,ggs.     `$$b:   Kernel: 5.10.0-11-amd64 
`d$$'     ,$P"'   .    $$$    Uptime: 37 days, 6 hours, 44 mins 
 $$P      d$'     ,    $$P    Packages: 1087 (dpkg) 
 $$:      $$.   -    ,d$$'    Shell: bash 5.1.4 
 $$;      Y$b._   _,d$P'      Resolution: 1024x768 
 Y$$.    `.`"Y$$$$P"'         Terminal: /dev/pts/0 
 `$$b      "-.__              CPU: Intel Xeon E5-2630 v4 (6) @ 2.199GHz 
  `Y$$                        GPU: 00:02.0 Vendor 1234 Device 1111 
   `Y$$.                      Memory: 8562MiB / 16008MiB 
     `$$b.
       `Y$$b.                                         
          `"Y$b._                                     
              `"""

root@core:~/hitokoto-go# 
```

### 开发机

```shell
nya@CandiFantasy:/mnt/d/Projects/go/hitokoto-go$ wrk -t8 -c1000 -d10s --latency http://localhost:8080
Running 10s test @ http://localhost:8080
  8 threads and 1000 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency    31.64ms    4.87ms 125.23ms   91.67%
    Req/Sec     3.97k   236.24     4.55k    81.06%
  Latency Distribution
     50%   30.34ms
     75%   32.33ms
     90%   35.36ms
     99%   52.81ms
  312527 requests in 10.01s, 235.11MB read
Requests/sec:  31207.84
Transfer/sec:     23.48MB
nya@CandiFantasy:/mnt/d/Projects/go/hitokoto-go$ neofetch
            .-/+oossssoo+/-.
        `:+ssssssssssssssssss+:`
      -+ssssssssssssssssssyyssss+-         nya@CandiFantasy 
    .ossssssssssssssssssdMMMNysssso.       ---------------- 
   /ssssssssssshdmmNNmmyNMMMMhssssss/      OS: Ubuntu 20.04.4 LTS on Windows 10 x86_64 
  +ssssssssshmydMMMMMMMNddddyssssssss+     Kernel: 5.10.102.1-microsoft-standard-WSL2 
 /sssssssshNMMMyhhyyyyhmNMMMNhssssssss/    Uptime: 9 hours, 4 mins 
.ssssssssdMMMNhsssssssssshNMMMdssssssss.   Packages: 800 (dpkg) 
+sssshhhyNMMNyssssssssssssyNMMMysssssss+   Shell: bash 5.0.17 
ossyNMMMNyMMhsssssssssssssshmmmhssssssso   Theme: Adwaita [GTK3] 
ossyNMMMNyMMhsssssssssssssshmmmhssssssso   Icons: Adwaita [GTK3] 
+sssshhhyNMMNyssssssssssssyNMMMysssssss+   Terminal: /dev/pts/1 
.ssssssssdMMMNhsssssssssshNMMMdssssssss.   CPU: AMD Ryzen 9 5900X (24) @ 3.693GHz 
 /sssssssshNMMMyhhyyyyhdNMMMNhssssssss/    GPU: b10e:00:00.0 Microsoft Corporation Device 008e 
  +sssssssssdmydMMMMMMMMddddyssssssss+     Memory: 2535MiB / 64282MiB 
   /ssssssssssshdmNNNNmyNMMMMhssssss/
    .ossssssssssssssssssdMMMNysssso.                               
      -+sssssssssssssssssyyyssss+-                                 
        `:+ssssssssssssssssss+:`
            .-/+oossssoo+/-.



nya@CandiFantasy:/mnt/d/Projects/go/hitokoto-go$

```

## 开发团队

- Nya Candy

## 非常感谢

- [一言开源社区](https://github.com/hitokoto-osc)
