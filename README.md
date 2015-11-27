# datahub-client

# Datahub-Client
----------------
### 开始

在安装GO(需要go1.4以上版本)语言和设置了[GOPATH](http://golang.org/doc/code.html#GOPATH)环境变量之后，安装datahub-client：

```shell
go get github.com/asiainfoLDP/datahub
```

启动datahub服务:
```shell
sudo $GOPATH/bin/datahub --daemon --token xxxxxxxxxxx
```

### Docker方式启动

```shell
docker build -t datahub .
docker run -d -e "DAEMON_TOKEN=xxxxxxxxxxx" -e "DAEMON_ENTRYPOINT=http://XXXXXXXXXX:YYYY" -p xxxx:35800 datahub
```

### 运行Datahub CLI

Datahub CLI是datahub-client的命令行客户端，用来执行datahub相关命令。

- dp        
    - Datapool管理
- repo      
    - Repository管理
- subs      
    - Subscrption管理
- login     
    - 登录到dataos.io
- pull      
    - 下载数据
- pub
    - 发布数据

### Datahub Client 命令行使用说明
---
#### NOTE：
- 如果没有额外说明，所有的命令在没有错误发生时，不在终端输出任何信息，只记录到日志中。错误信息会打印到终端。
- 所有的命令执行都会记录到日志中，日志级别分[TRACE] [INFO] [WARNNING] [ERROR] [FATAL]。
- 参数支持全名和简称两种形式，例如--type等同于-t。详情见命令帮助。
- 参数赋值支持空格和等号两种形式，例如--type=file等同于--type file。

#### 1. datapool相关命令

##### 1.1. 列出所有命令池

```shell
datahub dp
```
输出
```shell
{%DPNAME    %DPTYPE}
```
例子
```shell
$ datahub dp
dp1     regular file 
dp2     db2
dphere  hdfs
dpthere api
$
```

##### 1.2. 列出datapool详情

```shell
datahub dp $DPNAME
```
输出
```shell
%DPNAME %DPTYPE %DPCONN
{%REPO/%ITEM:%TAG       %LOCAL_TIME     %T}
```
例子
```shell
$ datahub dp dp1
dp1 regular file    /var/lib/datahub/dp1

repo1/item1:tag1        12:34 Oct 11 2015       pub
repo1/item1:tag2        15:00 Nov 2  2015       pub
repo1/item2:latest  10:00 Nov 1  2015       pull
cmcc/beijing:latest 10:00 Nov 1  2015       pull
$ 
```

##### 1.3. 创建数据池

- 目前只支持本地目录形式的数据池创建。使用参数dpconn指定绝对路径，当没有设定dpconn选项时，会默认创建到/var/lib/datahub。

```shell
datahub dp create $DPNAME [--type=$dptype]
[--conn=$dpconn]
```
输出
```
%msg
```
例子
```
$ datahub dp create dp1 --type=file --conn=/home/daemon/dp1
dp1 created as /home/daemon/dp1
$
```

##### 1.4. 删除数据池

- 删除数据池不会删除目标数据池已保存的数据。该dp有发布的数据项时，不能被删除。删除是在sqlite中标记状态，不真实删除。

```
datahub dp rm $DPNAME
输出
例子
$ datahub dp rm dp1
ok
$
```

#### 2. subs相关命令

##### 2.1. 列出所有已订阅项

```
datahub subs 
```
输出
```
{%REPO/%ITEM    %TYPE}
```
例子
```
$ datahub subs
cmcc/beijing        regular file
repo1/testing       api
$
```

##### 2.2. 列出已订阅item详情

```
datahub subs $REPO/$ITEM
```
输出
```
%REPO/%ITEM     %TYPE
DESCRIPTION:
%DESCRIPTION
METADATA:
%METADATA
{%ITEM:%TAGNAME %UPDATE_TIME    %INFO}
```
例子
```
$ datahub subs cmcc/beijing
cmcc/beijing    regular file
DESCRIPTION:
移动数据北京地区
METADATA:
BLABLABLA
beijing:chaoyang    15:34 Oct 12 2015       600M
beijing:daxing  16:40 Oct 13 2015       435M
beijing:shunyi  16:40 Oct 14 2015       324M
beijing:haidian 16:40 Oct 15 2015       988M
$
```

#### 3. pull命令

##### 3.1. 拉取某个item的tag


```
datahub pull $REPO/$ITEM:$TAG $DATAPOOL[://$LOCATION]
```
输出
```
%msg
```
例子
```
$ datahub pull cmcc/beijing:chaoyang dp1://cmccbj
OK.
$
```

#### 4. login命令

- login命令支持被动调用，用于datahub client与datahub server交互时作认证。并将认证信息保存到环境变量，免去后续指令重复输入认证信息。

##### 4.1. 登录到dataos.io

```
datahub login [--user=user]
```
输出
```
%msg
```
例子
```
$ datahub login
login: datahub
password: *******
[INFO]Authorization failed.
$
```

#### 5. pub相关命令
-pub分为发布一个DataItem和发布一个Tag。
- 发布DataItem必须指定DATAPOOL和DATAPOOL下的子路径LOCATION ,  可选参数 --accesstype, -t= 指定DataItem属性：public, private, 默认private
- 发布Tag必须指定TAGDETAIL , 用来指定Tag对应文件名，该文件必须存在于$DATAPOOL://$LOCATION内
- 可选参数--comment, -m= ,描述DataItem或者Tag

##### 5.1. 发布一个DataItem

```
datahub pub $REPOSITORY/$DATAITEM $DATAPOOL://$LOCATION --accesstype=public [private]  --comment=" "
```
输出
```
Pub success,  OK
```
例子
```
$./datahub pub music_1/migu mydp://dirmigu --accesstype=public --comment="migu music desc"
Pub success,  OK
```

##### 5.2.发布一个Tag

```
datahub pub $REPO/$ITEM:$Tag $TAGDETAIL --comment=" "
```
输出
```
Pub success, OK
```
例子
```
$ datahub pub music_1/migu:migu_user_info migu_user_info.txt
Pub success, OK
$
```
#### 6. help命令

- help提供datahub所有命令的帮助信息。

##### 6.1. 列出帮助

```
datahub help [$CMD] [$SUBCMD]
```
输出
```
Usage of %CMD %SUBCMD
{  %OPTION=%DEFAULT_VALUE     %OPTION_DESCRIPTION}
```
例子
```
$ datahub help dp create
Usage of dp create:

--conn=            datapool connection info
--name=            datapool name
--type, -T=file    datapool type
$ datahub help dp
Usage of dp:

dp [create | rm] <dpname>
```

