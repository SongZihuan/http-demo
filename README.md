# HTTP-Demo
基于Go+Gin框架的HTTP-Demo程序。

主要用于打印一些网络信息，便于测试。

## 命令行参数
### HTTP
Http默认监听端口`:3366`，可通过`--address`和`-http-address`参数更改。

### HTTPS
Https默认不监听，可通过`--https-address`参数开启监听。

HTTPS默认使用acme的DNS-01调整申请HTTPS证书，因此你需要配置：`--https-domain`域名，`--https-aliyun-dns-access-key`和`--https-aliyun-dns-access-secret`阿里云具有DNS权限的RAM用户的Key和Secret。
你还可以选择性配置`--https-email`，申请CA联系人你（不会体现在证书上），默认为`no-reply@example.com`。
你还可以选择性配置`--https-cert-dir`，保存证书和和账号信息（每个邮箱都会对应一个账户），默认为当前目录下`ssl-certs`文件夹。

## 环境变量
通过环境变量也可以设置参数，但是会被命令行参数覆盖

`DH_HTTP_ADDRESS`等价于`--http-address`。

`DH_HTTPS_ADDRESS`等价于`--https-address`。

`DH_HTTPS_DOMAIN`等价于`--https-domain`。

`DH_HTTPS_EMAIL`等价于`--https-email`。

`DH_HTTPS_CERT_DIR`等价于`--http-cert-dir`。

`DH_HTTPS_ALIYUN_KEY`等价于`--https-aliyun-dns-access-key`。

`DH_HTTPS_ALIYUN_SECRET`等价于`--https-aliyun-dns-access-secret`。

## 路由
`/` - 打印请求信息

`/ip` - 打印接收请求时对方的IP地址，未必为请求人的IP地址，可能是代理的地址。

`/client/ip` - 请取人的地址，通过请求头X-Forwarder-For等获取

`/timestamp` - 当前时间戳

`/datetime` - 当前时间

`/hello` - 打印欢迎信息

`/empty` - 返回204，无body

## 协议
本软件基于[MIT LICENSE](./LICENSE)协议发布。
