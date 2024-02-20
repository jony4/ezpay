# 艺爪付费 Go SDK

为网站App提供会员、订阅、内购付费功能。

- 使用地址：https://revenue.ezboti.com/dash/
- API文档地址：https://www.ezboti.com/docs/revenue/

## 使用方式

```go
client, err := NewClient(config)
if err != nil {
    log.Fatal(err)
}
info, err := client.CustomerInfo(ctx, &EZCustomerInfoRequest{})
if err != nil {
    log.Fatal(err)
}
```
