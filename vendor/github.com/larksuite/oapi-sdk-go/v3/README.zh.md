# 飞书开放接口SDK

旨在让开发者便捷的调用飞书开放API、处理订阅的消息事件、处理服务端推送的卡片行为。

## 目录


<!-- toc -->

- [安装](#安装)
- [API Client](#api-client)
    - [创建API Client](#创建api-client)
    - [配置API Client](#配置api-client)

- [API调用](#api调用)
    - [基本用法](#基本用法)
    - [设置请求选项](#设置请求选项)
    - [原生API调用方式](#原生api调用方式)

- [处理消息事件回调](#处理消息事件回调)
    - [基本用法](#基本用法-1)
    - [消息处理器内给对应租户发消息](#消息处理器内给对应租户发消息)
    - [集成gin框架](#集成gin框架)
        - [安装集成包](#安装集成包)
        - [集成示例](#集成示例)

- [处理卡片行为回调](#处理卡片行为回调)
    - [基本用法](#基本用法-2)
    - [返回卡片消息](#返回卡片消息)
    - [返回自定义消息](#返回自定义消息)
    - [卡片行为处理器内给对应租户发消息](#卡片行为处理器内给对应租户发消息)
    - [集成gin框架](#集成gin框架)
        - [安装集成包](#安装集成包)
        - [集成示例](#集成示例)

<!-- tocstop -->

## 安装

```go
go get -u github.com/larksuite/oapi-sdk-go/v3@v3.0.9
```

## API Client

开发者在调用 API 前，需要先创建一个 API Client，然后才可以基于 API Client 发起 API 调用。

### 创建API Client

- 对于自建应用，可使用下面代码来创建一个 API Client

```go
var client = lark.NewClient("appID", "appSecret") // 默认配置为自建应用
```

- 对于商店应用，需在创建 API Client 时，使用 lark.WithMarketplaceApp 方法指定 AppType 为商店应用

```go
var client = lark.NewClient("appID", "appSecret",lark.WithMarketplaceApp()) // 设置App为商店应用
```

### 配置API Client

创建 API Client 时，可对 API Client 进行一定的配置，比如我们可以在创建 API Client 时设置日志级别、设置 http 请求超时时间等等：

```go
var client = lark.NewClient("appID", "appSecret",
    lark.WithLogLevel(larkcore.LogLevelDebug),
    lark.WithReqTimeout(3*time.Second),
    lark.WithEnableTokenCache(true),
    lark.WithHelpdeskCredential("id", "token"),
    lark.WithHttpClient(http.DefaultClient))
```

每个配置选项的具体含义，如下表格：

<table>
  <thead align=left>
    <tr>
      <th>
        配置选项
      </th>
      <th>
        配置方式
      </th>
       <th>
        描述
      </th>
    </tr>
  </thead>
  <tbody align=left valign=top>
    <tr>
          <th>
            <code>AppType</code>
          </th>
          <td>
            <code>lark.WithMarketplaceApp()</code>
          </td>
          <td>
    设置 App 类型为 商店应用 ，ISV 开发者必须要设置该选项。
          </td>
    </tr>
    <tr>
      <th>
        <code>LogLevel</code>
      </th>
      <td>
        <code>lark.WithLogLevel(logLevel larkcore.LogLevel)</code>
      </td>
      <td>
设置 API Client 的日志输出级别(默认为 Info 级别)，枚举值如下：

- LogLevelDebug
- LogLevelInfo
- LogLevelWarn
- LogLevelError

</td>
</tr>

<tr>
      <th>
        <code>Logger</code>
      </th>
      <td>
        <code>lark.WithLogger(logger larkcore.Logger)</code>
      </td>
      <td>
设置API Client的日志器，默认日志输出到标准输出。

开发者可通过实现下面的 Logger 接口，来设置自定义的日志器:

```go
type Logger interface {
    Debug(context.Context, ...interface{})
    Info(context.Context, ...interface{})
    Warn(context.Context, ...interface{})
    Error(context.Context, ...interface{})
}
```

</td>
</tr>

<tr>
      <th>
        <code>LogReqAtDebug</code>
      </th>
      <td>
        <code>lark.WithLogReqAtDebug(printReqRespLog bool)</code>
      </td>
      <td>
设置是否开启 Http 请求参数和响应参数的日志打印开关；开启后，在 debug 模式下会打印 http 请求和响应的 headers,body 等信息。

在排查问题时，开启该选项，有利于问题的排查。

</td>
</tr>


<tr>
      <th>
        <code>BaseUrl</code>
      </th>
      <td>
        <code>lark.WithOpenBaseUrl(baseUrl string)</code>
      </td>
      <td>
设置飞书域名，默认为FeishuBaseUrl，可用域名列表为：

```go
// 飞书域名
var FeishuBaseUrl = "https://open.feishu.cn"

// Lark域名
var LarkBaseUrl = "https://open.larksuite.com"
```

</td>
</tr>

<tr>
      <th>
        <code>TokenCache</code>
      </th>
      <td>
        <code>lark.WithTokenCache(cache larkcore.Cache)</code>
      </td>
      <td>
设置 token 缓存器，用来缓存 token 和 appTicket, 默认实现为内存。

如开发者想要定制 token 缓存器，需实现下面 Cache 接口:

```go
type Cache interface {
  Set(ctx context.Context, key string, value string, expireTime time.Duration) error
  Get(ctx context.Context, key string) (string, error)
}
```

对于 ISV 开发者来说，如需要 SDK 来缓存 appTicket，需要实现该接口，实现提供分布式缓存。

</td>
</tr>


<tr>
      <th>
        <code>EnableTokenCache</code>
      </th>
      <td>
        <code>lark.WithEnableTokenCache(enableTokenCache bool)</code>
      </td>
      <td>
设置是否开启 TenantAccessToken 的自动获取与缓存。

默认开启，如需要关闭可传递 false。
</td>
</tr>

<tr>
      <th>
        <code>HelpDeskId、HelpDeskToken</code>
      </th>
      <td>
        <code>lark.WithHelpdeskCredential(helpdeskID, helpdeskToken string)</code>
      </td>
      <td>
该选项仅在调用服务台业务的 API 时需要配置。
</td>
</tr>


<tr>
      <th>
        <code>ReqTimeout</code>
      </th>
      <td>
        <code>lark.WithReqTimeout(time time.Duration)</code>
      </td>
      <td>
设置 SDK 内置的 Http Client 的请求超时时间，默认为0代表永不超时。
</td>
</tr>

<tr>
      <th>
        <code>HttpClient</code>
      </th>
      <td>
        <code>lark.WithHttpClient(httpClient larkcore.HttpClient)</code>
      </td>
      <td>
设置 HttpClient，用于替换 SDK 提供的默认实现。

开发者可通过实现下面的 HttpClient 接口来设置自定义的 HttpClient:

```go
type HttpClient interface {
  Do(*http.Request) (*http.Response, error)
}

```

</td>
</tr>

  </tbody>
</table>

## API调用
创建完毕 API Client，我们可以使用 ``Client.业务域.资源.方法名称`` 来定位具体的 API 方法，然后对具体的 API 发起调用。

![](doc/find_method.jpg)

飞书开放平台开放的所有 API 列表，可点击[这里查看](https://open.feishu.cn/document/ukTMukTMukTM/uYTM5UjL2ETO14iNxkTN/server-api-list)

### 基本用法

如下示例我们通过 client 调用文档业务的 Create 方法，创建一个文档：

``` go
import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/larksuite/oapi-sdk-go/v3"
	"github.com/larksuite/oapi-sdk-go/v3/core"
	"github.com/larksuite/oapi-sdk-go/v3/service/docx/v1"
)


func main() {
	// 创建 client
	client := lark.NewClient("appID", "appSecret")

	// 发起请求
	resp, err := client.Docx.Document.Create(context.Background(), larkdocx.NewCreateDocumentReqBuilder().
		Body(larkdocx.NewCreateDocumentReqBodyBuilder().
			FolderToken("token").
			Title("title").
			Build()).
		Build())

	//处理错误
	if err != nil {
           // 处理err
           return
	}

	// 服务端错误处理
	if !resp.Success() {
           fmt.Println(resp.Code, resp.Msg, resp.RequestId())
	   return 
	}

	// 业务数据处理
	fmt.Println(larkcore.Prettify(resp.Data))
}
```

更多 API 调用示例：[./sample/api/im.go](./sample/api/im.go)

### 设置请求选项

开发者在每次发起 API 调用时，可以设置请求级别的一些参数，比如传递 UserAccessToken ,自定义 Headers 等：

```go
import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/larksuite/oapi-sdk-go/v3"
	"github.com/larksuite/oapi-sdk-go/v3/core"
	"github.com/larksuite/oapi-sdk-go/v3/service/docx/v1"
)

func main() {
	// 创建client
	client := lark.NewClient("appID", "appSecret")

	// 自定义请求headers
	header := make(http.Header)
	header.Add("k1", "v1")
	header.Add("k2", "v2")

	// 发起请求
	resp, err := client.Docx.Document.Create(context.Background(), larkdocx.NewCreateDocumentReqBuilder().
		Body(larkdocx.NewCreateDocumentReqBodyBuilder().
			FolderToken("token").
			Title("title").
			Build(),
		).
		Build(),
		larkcore.WithHeaders(header), // 设置自定义headers
	)

	//处理错误
	if err != nil {
	   // 处理err
	   return
	}

	// 服务端错误处理
	if !resp.Success() {
	   fmt.Println(resp.Code, resp.Msg, resp.RequestId())
	   return
	}

	// 业务数据处理
	fmt.Println(larkcore.Prettify(resp.Data))
}

```

如下表格，展示了所有请求级别可设置的选项：

<table>
  <thead align=left>
    <tr>
      <th>
        配置选项
      </th>
      <th>
        配置方式
      </th>
       <th>
        描述
      </th>
    </tr>
  </thead>
  <tbody align=left valign=top>
    <tr>
      <th>
        <code>Header</code>
      </th>
      <td>
        <code>larkcore.WithHeaders(header http.Header)</code>
      </td>
      <td>
设置自定义请求头，开发者可在发起请求时，这些请求头会被透传到飞书开放平台服务端。

</td>
</tr>

<tr>
      <th>
        <code>UserAccessToken</code>
      </th>
      <td>
        <code>larkcore.WithUserAccessToken(userAccessToken string)</code>
      </td>
      <td>
设置用户token，当开发者需要以用户身份发起调用时，需要设置该选项的值。

</td>
</tr>

<tr>
      <th>
        <code>TenantAccessToken</code>
      </th>
      <td>
        <code>larkcore.WithTenantAccessToken(tenantAccessToken string)</code>
      </td>
      <td>
设置租户 token，当开发者自己维护租户 token 时（即创建Client时EnableTokenCache设置为了false），需通过该选项传递 租户 token。

</td>
</tr>

<tr>
      <th>
        <code>TenantKey</code>
      </th>
      <td>
        <code>larkcore.WithTenantKey(tenantKey string)</code>
      </td>
      <td>
设置租户 key, 当开发者开发商店应用时，必须设置该选项。
</td>
</tr>


<tr>
      <th>
        <code>RequestId</code>
      </th>
      <td>
        <code>larkCore.WithRequestId(requestId string)</code>
      </td>
      <td>
设置请求 ID，用来做请求的唯一标识，该 ID 会被透传到飞书开放平台服务端。

</td>
</tr>

  </tbody>
</table>

### 原生API调用方式

有些老版本的开放接口，不能生成结构化的 API， 导致 SDK 内无法提供结构化的使用方式，这时可使用原生模式进行调用：

```go
import (
	"context"
	"fmt"
	"os"

	"github.com/larksuite/oapi-sdk-go/v3"
	"github.com/larksuite/oapi-sdk-go/v3/core"
)

func main() {
	// 创建 API Client
	var appID, appSecret = os.Getenv("APP_ID"), os.Getenv("APP_SECRET")
	var cli = lark.NewClient(appID, appSecret, lark.WithLogReqAtDebug(true), lark.WithLogLevel(larkcore.LogLevelDebug))

	// 发起请求
	resp, err := cli.Do(context.Background(),
		&larkcore.ApiReq{
			HttpMethod:                http.MethodGet,
			ApiPath:                   "https://open.feishu.cn/open-apis/contact/v3/users/:user_id",
			Body:                      nil,
			QueryParams:               larkcore.QueryParams{"user_id_type": []string{"open_id"}},
			PathParams:                larkcore.PathParams{"user_id": "ou_c245b0a7dff2725cfa2fb104f8b48b9d"},
			SupportedAccessTokenTypes: []larkcore.AccessTokenType{larkcore.AccessTokenTypeUser},
		},
		larkcore.WithUserAccessToken("u-3Sr1oTO4V1FWxTFTFYuFCqhk2Vs4h5IbhMG00gmw0CXh"),
	)

	// 错误处理
	if err != nil {
		fmt.Println(err)
		return
	}

	// 获取请求 ID
	fmt.Println(resp.RequestId())

	// 处理请求结果
	fmt.Println(resp.StatusCode)      // http status code
	fmt.Println(resp.Header)          // http header
	fmt.Println(string(resp.RawBody)) // http body
}
```

更多 API 调用示例：[./sample/callrawapi/api.go](./sample/callrawapi/api.go)

## 处理消息事件回调
关于消息订阅相关的知识，可以点击[这里查看](https://open.feishu.cn/document/ukTMukTMukTM/uUTNz4SN1MjL1UzM)

飞书开放平台开放的所有事件列表，可点击[这里查看](https://open.feishu.cn/document/ukTMukTMukTM/uYDNxYjL2QTM24iN0EjN/event-list)
### 基本用法

开发者订阅消息事件后，可以使用下面代码，对飞书开放平台推送的消息事件进行处理，如下代码基于 go-sdk 原生 http server 启动一个 httpServer：

```go
import (
	"context"
	"fmt"
	"net/http"

	"github.com/larksuite/oapi-sdk-go/v3/core"
	"github.com/larksuite/oapi-sdk-go/v3/event"
	"github.com/larksuite/oapi-sdk-go/v3/event/dispatcher"
	"github.com/larksuite/oapi-sdk-go/v3/core/httpserverext"
	"github.com/larksuite/oapi-sdk-go/v3/service/contact/v3"
	"github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
)

func main() {
    // 注册消息处理器
    handler := dispatcher.NewEventDispatcher("verificationToken", "eventEncryptKey").OnP2MessageReceiveV1(func(ctx context.Context, event *larkim.P2MessageReceiveV1) error {
        // 处理消息 event，这里简单打印消息的内容 
        fmt.Println(larkcore.Prettify(event))
        fmt.Println(event.RequestId())
        return nil
    }).OnP2MessageReadV1(func(ctx context.Context, event *larkim.P2MessageReadV1) error {
        // 处理消息 event，这里简单打印消息的内容
        fmt.Println(larkcore.Prettify(event))
        fmt.Println(event.RequestId())
        return nil
    })
    
    // 注册 http 路由
    http.HandleFunc("/webhook/event", httpserverext.NewEventHandlerFunc(handler, larkevent.WithLogLevel(larkcore.LogLevelDebug)))
    
    // 启动 http 服务
    err := http.ListenAndServe(":9999", nil)
    if err != nil {
        panic(err)
    }
}


```

其中 NewEventDispatcher 方法的参数用于签名验证和消息解密使用，默认可以传递为空串；但是如果开发者的应用在 [控制台](https://open.feishu.cn/app?lang=zh-CN) 的【事件订阅】里面开启了加密，则必须传递控制台上提供的值。

![Console](doc/console.jpeg)

需要注意的是注册处理器时，比如使用 OnP2MessageReceiveV1 注册接受消息事件回调时，其中的P2为消息协议版本，当前飞书开放平台存在 [两种消息协议](https://open.feishu.cn/document/ukTMukTMukTM/uUTNz4SN1MjL1UzM#8f960a4b) ，分别为1.0和2.0。

如下图开发者在注册消息处理器时，需从 [事件列表](https://open.feishu.cn/document/ukTMukTMukTM/uYDNxYjL2QTM24iN0EjN/event-list) 中查看自己需要的是哪种协议的事件。
如果是1.0的消息协议，则注册处理器时，需要找以OnP1xxxx开头的。如果是2.0的消息协议，则注册处理器时，需要找以OnP2xxxx开头的。




![Console](doc/event_protocol.jpeg)

更多事件订阅示例：[./sample/event/event.go](./sample/event/event.go)

## 消息处理器内给对应租户发消息
针对 ISV 开发者，如果想在消息处理器内给对应租户的用户发送消息，则需先从消息事件内获取租户 key,然后使用下面方式调用消息 API 进行消息发送：

```go
import (
	"context"
	"fmt"
	"net/http"

	"github.com/larksuite/oapi-sdk-go/v3/core"
	"github.com/larksuite/oapi-sdk-go/v3/event"
	"github.com/larksuite/oapi-sdk-go/v3/event/dispatcher"
	"github.com/larksuite/oapi-sdk-go/v3/core/httpserverext"
	"github.com/larksuite/oapi-sdk-go/v3/service/contact/v3"
	"github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
)

func main() {
    // 注册消息处理器
    handler := dispatcher.NewEventDispatcher("verificationToken", "eventEncryptKey").OnP2MessageReceiveV1(func(ctx context.Context, event *larkim.P2MessageReceiveV1) error {
        // 处理消息 event，这里简单打印消息的内容 
        fmt.Println(larkcore.Prettify(event))
        fmt.Println(event.RequestId())
        
        // 获取租户 key 并发送消息
        tenanKey := event.TenantKey()
        
        // ISV 给指定租户发送消息
        resp, err := client.Im.Message.Create(context.Background(), larkim.NewCreateMessageReqBuilder().
                ReceiveIdType(larkim.ReceiveIdTypeOpenId).
                Body(larkim.NewCreateMessageReqBodyBuilder().
                    MsgType(larkim.MsgTypePost).
                    ReceiveId("ou_c245b0a7dff2725cfa2fb104f8b48b9d").
                    Content("text").
                    Build(), larkcore.WithTenantKey(tenanKey)).
                Build())
                
        // 发送结果处理，resp,err
		
        return nil
    })
    
    // 注册 http 路由
    http.HandleFunc("/webhook/event", httpserverext.NewEventHandlerFunc(handler, larkevent.WithLogLevel(larkcore.LogLevelDebug)))
    
    // 启动 http 服务
    err := http.ListenAndServe(":9999", nil)
    if err != nil {
        panic(err)
    }
}

```


更多事件订阅示例：[./sample/event/event.go](./sample/event/event.go)


### 集成Gin框架
如果开发者当前应用使用的是 Gin Web 框架，并且不想要使用 Go-Sdk 提供的原生的 Http Server，则可使用下面方式，把当前应用的 Gin 服务与 SDK进行集成。

要想把 SDK 集成已有 Gin 框架，开发者需要引入集成包 [oapi-sdk-gin](https://github.com/larksuite/oapi-sdk-gin)

#### 安装集成包

```go
go get -u github.com/larksuite/oapi-sdk-gin
```

#### 集成示例

```go
import (
	"context"
	"fmt"

	 "github.com/gin-gonic/gin"
	 "github.com/larksuite/oapi-sdk-gin"
	 "github.com/larksuite/oapi-sdk-go/v3/card"
	 "github.com/larksuite/oapi-sdk-go/v3/core"
	 "github.com/larksuite/oapi-sdk-go/v3/event/dispatcher"
	 "github.com/larksuite/oapi-sdk-go/v3/service/contact/v3"
	 "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
)

func main() {
	// 注册消息处理器
	handler := dispatcher.NewEventDispatcher("verificationToken", "eventEncryptKey").OnP2MessageReceiveV1(func(ctx context.Context, event *larkim.P2MessageReceiveV1) error {
		fmt.Println(larkcore.Prettify(event))
		fmt.Println(event.RequestId())
		return nil
	}).OnP2MessageReadV1(func(ctx context.Context, event *larkim.P2MessageReadV1) error {
		fmt.Println(larkcore.Prettify(event))
		fmt.Println(event.RequestId())
		return nil
	}).OnP2UserCreatedV3(func(ctx context.Context, event *larkcontact.P2UserCreatedV3) error {
		fmt.Println(larkcore.Prettify(event))
		fmt.Println(event.RequestId())
		return nil
	})

	...

	// 在已有 Gin 实例上注册消息处理路由
	gin.POST("/webhook/event", sdkginext.NewEventHandlerFunc(handler))
}
```


## 处理卡片行为回调

关于卡片行为相关的知识，可点击[这里查看](https://open.feishu.cn/document/ukTMukTMukTM/uczM3QjL3MzN04yNzcDN)
### 基本用法

开发者配置消息卡片回调地址后，可以使用下面代码，对飞书开放平台推送的卡片行为进行处理，如下代码基于go-sdk原生http server启动一个httpServer：

```go
import (
	"context"
	"fmt"
	"net/http"

	"github.com/larksuite/oapi-sdk-go/v3/card"
	"github.com/larksuite/oapi-sdk-go/v3/core"
	"github.com/larksuite/oapi-sdk-go/v3/core/httpserverext"
)

func main() {
	// 创建 card 处理器
	cardHandler := larkcard.NewCardActionHandler("v", "", func(ctx context.Context, cardAction *larkcard.CardAction) (interface{}, error) {
		// 处理 cardAction, 这里简单打印卡片内容
		fmt.Println(larkcore.Prettify(cardAction))
	    fmt.Println(cardAction.RequestId())
		// 无返回值示例
		return nil, nil
	})

	// 注册处理器
	http.HandleFunc("/webhook/card", httpserverext.NewCardActionHandlerFunc(cardHandler, larkevent.WithLogLevel(larkcore.LogLevelDebug)))

	// 启动 http 服务
	err := http.ListenAndServe(":9999", nil)
	if err != nil {
		panic(err)
	}
}

```

如上示例，如果不需要处理器内返回业务结果给飞书服务端，则直接使用这种无返回值用法

更多卡片行为处理示例：[./sample/card/card.go](./sample/card/card.go)

### 返回卡片消息

如开发者需要卡片处理器内同步返回用于更新消息卡片的消息体，则可使用下面方法方式进行处理：

```go

import (
	"context"
	"fmt"
	"net/http"

	"github.com/larksuite/oapi-sdk-go/v3/card"
	"github.com/larksuite/oapi-sdk-go/v3/core"
	"github.com/larksuite/oapi-sdk-go/v3/core/httpserverext"
)

func main() {
	// 创建card处理器
	cardHandler := larkcard.NewCardActionHandler("v", "", func(ctx context.Context, cardAction *larkcard.CardAction) (interface{}, error) {
		fmt.Println(larkcore.Prettify(cardAction))
	    fmt.Println(cardAction.RequestId())
		
		// 创建卡片信息
		messageCard := larkcard.NewMessageCard().
		Config(config).
		Header(header).
		Elements([]larkcard.MessageCardElement{divElement, processPersonElement}).
		CardLink(cardLink).
		Build()

		return messageCard, nil
	})

	// 注册处理器
	http.HandleFunc("/webhook/card", httpserverext.NewCardActionHandlerFunc(cardHandler, larkevent.WithLogLevel(larkcore.LogLevelDebug)))

	// 启动http服务
	err := http.ListenAndServe(":9999", nil)
	if err != nil {
		panic(err)
	}
}

```

更多卡片行为处理示例：[./sample/card/card.go](./sample/card/card.go)

### 返回自定义消息

如开发者需卡片处理器内返回自定义内容，则可以使用下面方式进行处理：

```go 
import (
	"context"
	"fmt"
	"net/http"

	"github.com/larksuite/oapi-sdk-go/v3/card"
	"github.com/larksuite/oapi-sdk-go/v3/core"
	"github.com/larksuite/oapi-sdk-go/v3/core/httpserverext"
)

func main() {
	// 创建 card 处理器
	cardHandler := larkcard.NewCardActionHandler("v", "", func(ctx context.Context, cardAction *larkcard.CardAction) (interface{}, error) {
		fmt.Println(larkcore.Prettify(cardAction))
	    fmt.Println(cardAction.RequestId())
		
		// 创建 http body
		body := make(map[string]interface{})
		body["content"] = "hello"

		i18n := make(map[string]string)
		i18n["zh_cn"] = "你好"
		i18n["en_us"] = "hello"
		i18n["ja_jp"] = "こんにちは"
		body["i18n"] = i18n 
		
		// 创建自定义消息：http状态码，body内容
		resp := &larkcard.CustomResp{
			StatusCode: 400,
			Body:       body,
		}

		return resp, nil
	})

	// 注册处理器
	http.HandleFunc("/webhook/card", httpserverext.NewCardActionHandlerFunc(cardHandler, larkevent.WithLogLevel(larkcore.LogLevelDebug)))

	// 启动 http 服务
	err := http.ListenAndServe(":9999", nil)
	if err != nil {
		panic(err)
	}
}

```

更多卡片行为处理示例：[./sample/card/card.go](./sample/card/card.go)


### 卡片行为处理器内给对应租户发消息

针对 ISV 开发者，如果想在卡片行为处理器内给对应租户的用户发送消息，则需先从卡片行为内获取租户 key ,然后使用下面方式调用消息 API 进行消息发送：


```go
import (
	"context"
	"fmt"
	"net/http"

	"github.com/larksuite/oapi-sdk-go/v3/card"
	"github.com/larksuite/oapi-sdk-go/v3/core"
	"github.com/larksuite/oapi-sdk-go/v3/core/httpserverext"
)

func main() {
	// 创建 card 处理器
	cardHandler := larkcard.NewCardActionHandler("v", "", func(ctx context.Context, cardAction *larkcard.CardAction) (interface{}, error) {
        
        // 处理 cardAction, 这里简单打印卡片内容  
        fmt.Println(larkcore.Prettify(cardAction))
        fmt.Println(cardAction.RequestId())
	    
        // 获取租户 key 并发送消息
        tenanKey := cardAction.TenantKey
        
        // ISV 给指定租户发送消息
        resp, err := client.Im.Message.Create(context.Background(), larkim.NewCreateMessageReqBuilder().
                ReceiveIdType(larkim.ReceiveIdTypeOpenId).
                Body(larkim.NewCreateMessageReqBodyBuilder().
                    MsgType(larkim.MsgTypePost).
                    ReceiveId("ou_c245b0a7dff2725cfa2fb104f8b48b9d").
                    Content("text").
                    Build(), larkcore.WithTenantKey(tenanKey)).
                Build())
                
        // 发送结果处理，resp,err
		
        return nil, nil
	})

	// 注册处理器
	http.HandleFunc("/webhook/card", httpserverext.NewCardActionHandlerFunc(cardHandler, larkevent.WithLogLevel(larkcore.LogLevelDebug)))

	// 启动 http 服务
	err := http.ListenAndServe(":9999", nil)
	if err != nil {
		panic(err)
	}
}

```

更多卡片行为处理示例：[./sample/card/card.go](./sample/card/card.go)


### 集成gin框架

如果开发者当前应用使用的是 Gin Web 框架，并且不想要使用 Go-Sdk 提供的原生的 Http Server，则可使用下面方式，把当前应用的 Gin 服务与 SDK进行集成。

要想把 SDK 集成已有 Gin 框架，开发者需要引入集成包 [oapi-sdk-gin](https://github.com/larksuite/oapi-sdk-gin)

#### 安装集成包

```go
go get -u github.com/larksuite/oapi-sdk-gin
```

#### 集成示例

```go
import (
    "context"
    "fmt"
    
    "github.com/gin-gonic/gin"
    "github.com/larksuite/oapi-sdk-gin"
    "github.com/larksuite/oapi-sdk-go/v3/card"
    "github.com/larksuite/oapi-sdk-go/v3/core"
)


func main() {
      // 创建 card 处理器
      cardHandler := larkcard.NewCardActionHandler("v", "", func(ctx context.Context, cardAction *larkcard.CardAction) (interface{}, error) {
      fmt.Println(larkcore.Prettify(cardAction))
      fmt.Println(cardAction.RequestId())
    
      return nil, nil
      })
      ...
      // 在已有的 Gin 实例上注册卡片处理路由
      gin.POST("/webhook/card", sdkginext.NewCardActionHandlerFunc(cardHandler))
      ...
}
```

## License
使用 MIT



