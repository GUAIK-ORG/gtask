# GTask-任务系统

## 简介

在某些场景下我们可能需要定时去处理一些任务，我们需要编写服务端的代码让它周期性的去工作，例如检查客户端的在线状态、定期的去获取接口的数据等，在一个系统中可能会存在多种不同类型的定时任务需要开发者去开发和维护。

GTask支持远程管理任务，执行任务脚本（目前仅支持lua），开发者可以创建任务实例并为实例添加处理器，在处理器中配置执行脚本就可以轻松的管理任务。

## 编译与使用

### Linux

```bash
./build/linux.sh
```

### MaxOS

```bash
./build/macos.sh
```

编译成功后在release目录下会生成两个文件，`gtask`和`client`，分别是服务端程序和客户端程序。

### 服务端运行

在`config/config.json`中保存着服务端监听的端口号和登陆的密钥，在客户端连接服务器的时候将会用到该密钥。启动服务端时需要指定配置参数。

```bash
./release/gtask -cfg ./config/config.json
```

### 客户端连接服务器

```bash
./release/client -h [host] -p 1126
```

这时会要求输入密钥，该密钥是初始化服务器时候的密钥。

## 客户端指令说明介绍

|  指令   | 参数  | 参数介绍 | 示例 |
|  ----  | ----  | ---- | ---- |
| create job  | key | 任务key | create job test |
| use  | key | 该key必须已经被创建 | use test |
| create processor  | [filePath trigger bReset bLoop bExit] | （必须先use job） [脚本文件 触发时间(秒) 能否被重置 是否循环 是否退出] | create processor ./example/lua/test_json.lua 3 0 1 0 |
| run  |  | （必须先use job）启动任务 | run |
| delete  |  | （必须先use job）停止并删除任务 | delete |

### LUA支持函数

### 加密

|  函数   | 参数 | 返回值 | 描述 |
|  ----  | ---- | ---- | ---- |
| md5 | string | string | 获取字符串的MD5值 |
| base64 | string | string | 获取字符串的base64编码 |
| base64UrlSafe | string | string | 获取字符串url安全的base64编码 |
| base64 | string | string | 获取字符串的base64编码 |
| hmac | [key:string, str:string] | string | 获取hmac值 |
| sha1 | string | string | 获取字符串的SHA1值 |

### 时间

|  函数   | 参数 | 返回值 | 描述 |
|  ----  | ---- | ---- | ---- |
| now |  | number | 获取当前时间戳（毫秒） |

### Json

|  函数   | 参数 | 返回值 | 描述 |
|  ----  | ---- | ---- | ---- |
| jsonMarshal | table | string | 将table转换成json字符串 |
| jsonUnMarshal | string | table | 将json字符串转转换成table |

### Http

|  函数   | 参数 | 返回值 | 描述 |
|  ----  | ---- | ---- | ---- |
| httpGet | [url:string header:table] | [res:string ok:bool] | 发送GET请求 |
| httpPost | [url:string header:table body:string] | [res:string ok:bool] | 发送POST请求 |

## 测试

使用默认配置文件并将服务端运行在本机（localhost）

```bash
./release/client -h localhost -p 1126
secretKey:647851f2fcf6101aefa4a2c59a329a11c60300a4

# 创建任务
> create job test
create job [test] success

# 选择任务
> use test
select job [test]

# 为任务创建执行器,解析json并打印相关数据，3秒执行一次，循环执行
test > create processor ./example/lua/test_json.lua 3 0 1 0
create processor success

# 运行任务
test > run
run job [test] success

# 为任务补充一条执行器，用来发送GET请求获取网站数据，5秒执行一次，循环执行
test > create processor ./example/lua/test_http_get.lua 5 0 1 0
create processor success

# 停止并删除当前任务
test > delete
>
```

## LUA执行代码介绍

```lua
// 函数名必须是processor
// key(string): 任务key
// count(number): 时间计数
function processor(key,count)
    data = {}
    data["hello"]="world"
    data["a"] = {}
    data["a"]["b"] = "b"
    data["a"]["c"] = {1,2,3,4,5,6}
    res = jsonMarshal(data)
    res = jsonUnMarshal(res)
    for k,v in ipairs(res["a"]["c"]) do
        print(k,v)
    end
    return true　// 如果返回false，当前执行器将会退出
end
```
