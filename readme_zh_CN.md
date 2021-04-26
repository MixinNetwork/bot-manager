
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

#### 1. 环境变量的路径
> 1. 开发环境： `/client/.env` 和 `/client/.env.development`
> 2. 生产环境： `/client/.env.production`

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
