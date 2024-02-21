# 艺爪付费 Go SDK

为网站App提供会员、订阅、内购付费功能。

- 使用地址：https://revenue.ezboti.com/dash/
- API文档地址：https://www.ezboti.com/docs/revenue/

## 使用方式

```go
func main() {
    config := &ezpay.Config{
        ProjectID:     "1",
        ProjectSecret: "2",
        PaywallID:     "3",
    }
    client, err := ezpay.NewClient(config)
    if err != nil {
        log.Fatal(err)
    }
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    customerInfo, err := client.CustomerInfo(ctx, &ezpay.EZCustomerInfoRequest{})
    if err != nil {
        log.Fatal(err)
    }
    log.Info(customerInfo)
}
```

## 艺爪付费如果新增了接口如何快速拓展

- 参考`ezpay`目录下的`customer.info.go`文件，增加必须得请求与响应结构体。
- 调用 `client.Call()` 方法，传入请求结构体，解析响应内容。

示例：
```go
	response, err := c.Call(ctx, APICustomerInfo, request)
	if err != nil {
		return nil, err
	}
```