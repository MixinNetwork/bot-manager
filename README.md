### 1. Server-side deployment

#### 1. Copy the configuration file to the specified file name `app.conf`
```shell script
cd conf
cp app.conf.default app.conf
```

#### 2. Fill in the configuration file
```shell script
appname = bot-manager # the name of the app
httpport = 9098 # http service running port
runmode = dev # Load the specified environment variable
autorender = false
copyrequestbody = true

[dev] # Local development configuration
# Database please fill in the corresponding information of postgresql
dbHost=localhost
dbUser=
dbName=
dbPass=

# Please fill in the key-store of the mixin robot
clientId=
clientSecret=
sessionId=
privateKey=

# Customize the salt value of the salting algorithm
claimKey=dev


[prod] # Production environment configuration
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

#### 3. Run the development environment

This project is developed based on `beego`,
For installation of Bee tool, please refer to this link.
[https://beego.me/docs/install/](https://beego.me/docs/install/)

```shell script
bee run
```

#### 4. Production environment

1. In the production environment directory, create a new folder and file `conf/app.conf` and fill in the information,
> Note, `runmode=prod`
```shell script
go build
./bot-manager
```

### 2. Client deployment

#### 1. The path of the environment variable
> 1. Development environment: `/client/.env` and `/client/.env.development`
> 2. Production environment: `/client/.env.production`

#### 2. Configuration of environment variables

Example:
```shell script
VUE_APP_SERVER=http://localhost:9098/v1 # http service domain name
VUE_APP_WS_SERVER=ws://localhost:9099 # domain name of wss service
VUE_APP_CLIENT_ID= # client_id eg. 11efbb75-e7fe-44d7-a14f-698535289310
VUE_APP_SCOPE= # Need to request permissions, eg. SNAPSHOTS:READ+PROFILE:READ+ASSETS:READ
```

#### 3. Run the development environment

Running this project requires node environment support,
For the installation of Node Cemetery, please refer to this link:
[http://nodejs.cn/download/](http://nodejs.cn/download/)

```shell script
npm install
npm run serve
```

#### 4. Production environment
```shell script
npm install
npm run build
```