# Resty 客户端根证书方法

Source: https://github.com/go-resty/docs/blob/main/content/docs/root-certificates.md

提供了 Resty 客户端上可用的管理根证书的方法的概述，包括从文件、带监视器的文件和字符串设置。

```APIDOC
Client.SetRootCertificates(paths ...string)
  - 从文件路径添加一个或多个 PEM 编码的根证书。
  - 参数：
    - paths: PEM 编码证书文件的文件路径的可变参数列表。
  - 返回值：Resty 客户端实例，用于链式调用。

Client.SetRootCertificatesWatcher(opts *CertWatcherOptions, paths ...string)
  - 从文件路径添加一个或多个 PEM 编码的根证书，并带有用于动态重新加载的监视器。
  - 参数：
    - opts: 证书监视器的配置选项（例如，PoolInterval）。
    - paths: PEM 编码证书文件的文件路径的可变参数列表。
  - 返回值：Resty 客户端实例，用于链式调用。

Client.SetRootCertificateFromString(cert string)
  - 从字符串添加单个 PEM 编码的根证书。
  - 参数：
    - cert: 包含 PEM 编码证书的字符串。
  - 返回值：Resty 客户端实例，用于链式调用。

CertWatcherOptions struct {
  PoolInterval time.Duration // 检查证书修改的间隔。默认为 24 小时。
}
```

--------------------------------
