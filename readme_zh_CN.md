### 项目介绍
这是一个比较方便的机器人消息管理以及用户管理的后台系统。您可以比较方便的处理其他人向机器人发送的消息。

现已支持如下功能：
> 1. 后台直接回复用户消息。
> 2. 管理员通过 Messenger 客户端通过回复消息的方式，直接回复给用户。
> 3. 拒绝接收指定用户的信息。
> 4. 通过机器人发送公告以及公告的撤回
> 5. 指定关键字自动回复(被加好友自动回复)
> 6. 多管理员对同一机器人的支持(数据同步)
> 7. 统计机器人每日的新增用户、新增留言次数以及总用户数、总留言次数。
> 8. 设置关联分享的应用。

以下为应用内的部分截图。
![image](https://github.com/MixinNetwork/bot-manager/blob/main/img/data.png)
![image](https://github.com/MixinNetwork/bot-manager/blob/main/img/user.png)
![image](https://github.com/MixinNetwork/bot-manager/blob/main/img/message.png)
![image](https://github.com/MixinNetwork/bot-manager/blob/main/img/setting.png)



### 1. 服务端部署

#### 1. 复制配置文件为指定文件名 `app.conf`
```shell script
cd conf
cp app.conf.default app.conf
```

#### 2. 填写配置文件
```shell script
appname = bot-manager  # app 的名称
httpport = 9098 # http 服务运行的端口
runmode = dev # 加载的指定的环境变量
autorender = false
copyrequestbody = true

[dev] # 本地开发配置
# 数据库请填写 postgresql 相应信息
dbHost=localhost  
dbUser= 
dbName=
dbPass=

# 请填写 mixin 机器人的 key-store
clientId=
clientSecret=
sessionId=
privateKey=

# 自定义加盐算法的盐值
claimKey=dev


[prod] # 生产环境配置
dbHost=
dbUser=
dbName=
dbPass=

clientId=
clientSecret=
sessionId=
privateKey=

claimKey=prod
```

#### 3. 运行开发环境

本项目基于 `beego` 开发，
Bee工具的安装，请参照此链接。
[https://beego.me/docs/install/](https://beego.me/docs/install/) 

```shell script
bee run
```

#### 4. 生产环境

1. 在生产环境目录下，新建文件夹及文件 `conf/app.conf` 并填写好信息，
> 注意，`runmode=prod`
```shell script
go build 
./bot-manager
```

### 2. 客户端部署

#### 1. 复制配置文件为指定文件名
> 1. 开发环境： `/client/.env` 和 `/client/.env.development`
```shell script
cp /client/.env.default /client/.env.development
```
> 2. 生产环境： `/client/.env.production`
```shell script
cp /client/.env.default /client/.env.development
```

#### 2. 环境变量的配置

示例：
```shell script
VUE_APP_SERVER=http://localhost:9098/v1 # http服务的域名
VUE_APP_WS_SERVER=ws://localhost:9099 # wss服务的域名
VUE_APP_CLIENT_ID= # client_id eg. 11efbb75-e7fe-44d7-a14f-698535289310
VUE_APP_SCOPE= # 需要请求的权限， eg. SNAPSHOTS:READ+PROFILE:READ+ASSETS:READ
```

#### 3. 运行开发环境

运行本项目需要 node 环境支持， 
Node公墓的安装，请参照此链接：
[http://nodejs.cn/download/](http://nodejs.cn/download/)

```shell script
npm install
npm run serve
```

#### 4. 生产环境
```shell script
npm install
npm run build
```
