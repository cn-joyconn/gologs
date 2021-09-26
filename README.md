## 日志使用说明
    可以配置多个日志实体，每个实体可以指定日志输出类型、输出序列化方式。
    调用方式 
        import (log "github.com/cn-joyconn/gologs")
        log.Logger("myLogger").Info("我是日志")
## 日志配置文件说明
配置文件位置 `` ./conf/log.yml ``
### 示例
```
logs: #配置项根目录
  - {
    name: default,
    adapter : console,
    formatter: '%w %t | %F:%n>> %m',
    conf: {
      level: 1,
    }
  }    
  - {
    name: myLogger,
    adapter : file,
    formatter: '%w %t | %F:%n>> %m' ,
    conf: {
      filename: ./logs/myLogger.log, #保存的文件名
      maxsize: 100, #每个文件保存的最大尺寸，默认值是 1 << 28, //256 MB
      level: 1 , #日志保存的时候的级别，默认是 Trace 级别
      compress: false #压缩存储
    }     
  }
```
### Name
 logger名称

### Adapter
logger适配器，支持(console,file,multifile,conn,smtp,ElasticSearch)

### Adapter
日志输出格式
```
例：'%w %t | %F:%n>> %m'  
    2021-01-01 00:00:00.231 INFO | /data/beego-xadmin/test 21 我是日志内容

```
支持的占位符 ：  
$\color{#FF0000}{w}$ : 时间  
$\color{#FF0000}{m}$ : 消息  
$\color{#FF0000}{f}$ : 文件名  
$\color{#FF0000}{F}$ : 文件全路径  
$\color{#FF0000}{n}$ : 行数  
$\color{#FF0000}{l}$ : 消息级别，数字表示   
$\color{#FF0000}{t}$ : 消息级别，简写，例如`[I]`代表 INFO  
$\color{#FF0000}{T}$ : 消息级别，全称   


### Conf
日志对应的输出配置
```
输出到命令行，默认输出到`os.Stdout`：
Adapter = console
#     - level 输出的日志级别
#     - color 是否开启打印日志彩色打印(需环境支持彩色输出)

输出到文件
Adapter = file
#     - filename 保存的文件名
#     - maxsize 每个文件保存的最大尺寸，默认值是 1 << 28, //256 MB
#     - daily 是否按照每天 logrotate，默认是 true
#     - maxdays 文件最多保存多少天，默认保存 7 天
#     - rotate 是否开启 logrotate，默认是 true
#     - level 日志保存的时候的级别，默认是 Trace 级别
#     - perm 日志文件权限

输出到文件(需要单独写入文件的日志级别)
Adapter = multifile
#     - filename 保存的文件名
#     - maxsize 每个文件保存的最大尺寸，默认值是 1 << 28, //256 MB
#     - daily 是否按照每天 logrotate，默认是 true
#     - maxdays 文件最多保存多少天，默认保存 7 天
#     - rotate 是否开启 logrotate，默认是 true
#     - level 日志保存的时候的级别，默认是 Trace 级别
#     - perm 日志文件权限
#     - separate 需要单独写入文件的日志级别,设置后命名类似 test.error.log

网络输出
Adapter = conn
#     - reconnectOnMsg 是否每次链接都重新打开链接，默认是 false
#     - reconnect 是否自动重新链接地址，默认是 false
#     - net 发开网络链接的方式，可以使用 tcp、unix、udp 等
#     - addr 网络链接的地址
#     - level  日志保存的时候的级别，默认是 Trace 级别

邮件发送
Adapter = smtp
#     - username smtp 验证的用户名
#     - password smtp 验证密码
#     - host  发送的邮箱地址
#     - sendTos   邮件需要发送的人，支持多个
#     - subject   发送邮件的标题，默认是 `Diagnostic message from server`
#     - level 日志发送的级别，默认是 Trace 级别

输出到 ElasticSearch
Adapter = ElasticSearch
```