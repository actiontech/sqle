# go-workwx

![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/xen0n/go-workwx)
![GitHub tag (latest by date)](https://img.shields.io/github/v/tag/xen0n/go-workwx)
![GitHub Workflow Status (develop branch)](https://img.shields.io/github/actions/workflow/status/xen0n/go-workwx/go.yml?branch=develop)
[![GitHub license info](https://img.shields.io/github/license/xen0n/go-workwx)](./LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/xen0n/go-workwx)](https://goreportcard.com/report/github.com/xen0n/go-workwx)
[![Go Reference](https://pkg.go.dev/badge/github.com/xen0n/go-workwx.svg)](https://pkg.go.dev/github.com/xen0n/go-workwx)
![GitHub Repo stars](https://img.shields.io/github/stars/xen0n/go-workwx?style=social)

[**本项目需要您的帮助！**][maintainer-needed]

[maintainer-needed]: https://github.com/xen0n/go-workwx/issues/73

```go
import (
    "github.com/xen0n/go-workwx" // package workwx
)
```

A Work Weixin (a.k.a.  Wechat Work) client SDK for Golang, that happens to be
battle-tested and pretty serious about its types.

In production since late 2018, pushing all kinds of notifications and alerts
in at least 2 of Qiniu's internal systems.

一个 Golang 企业微信客户端 SDK；碰巧在生产环境试炼过，还对类型很严肃。

自 2018 年末以来，在七牛至少 2 个内部系统运转至今，推送各种通知、告警。


> English translation TODO for now, as the service covered here is not available
> outside of China (AFAIK).


## Why another wheel?

工作中需要用 Go 实现一个简单的消息推送，想着找个开源库算了，然而现有唯一的开源企业微信 Golang SDK 代码质量不佳。只好自己写一个。

*Update*: 自从这个库写出来，已经过了很久；现在（2019.08）已经有三四个同类项目了。
不过看了看这些“竞品”，发现自己这个库的类型设计、公开接口、access token 处理等方面还不赖。
为什么人们总是喜欢写死请求 `Host`、用全局量、为拆包而拆包甚至不惜公开内部方法呢？


## Supported Go versions

本库的 v1.x 分支最低支持到 `go1.17`。

CI 会在 `go1.17` 和 Go 的当前稳定版本、上一个稳定版本上跑测试，只有测试全部通过才可能合并 PR。

## Features

* 包名短
* 支持覆盖 API `Host`，用于自己拦一层网关、临时调试等等奇葩需求
* 支持使用自定义 `http.Client`
* access token 处理靠谱
    - 你可以直接就做 API 调用，会自动请求 access token
    - 你也可以一行代码起一个后台 access token 刷新 goroutine
    - 自带指数退避重试
* 严肃对待类型、公开接口
    - 公开暴露接口最小化，两步构造出 `WorkwxApp` 对象，然后直接用
    - 刻意不暴露企业微信原始接口请求、响应类型
    - 后续可能会做一个 `lowlevel` 包暴露裸的 API 接口，但很可能不做
    - 不为多态而多态，宁可 SDK 内部重复代码，也保证一个接口一类动作，下游用户 static dispatch
    - 个别数据模型做了调整甚至重做（如 `UserInfo`、`Recipient`），以鼓励 idiomatic Go 风格
    - *几乎*不会越俎代庖，一言不合 `panic`。**现存的少数一些情况都是要修掉的。**
* 自带一个 `workwxctl` 命令行小工具帮助调试
    - 用起来不爽提 issue 让我知道你在想啥

详情看 godoc 文档，还提供 Examples 小段代码可以参考。


## Supported APIs

* [x] 通讯录管理 (**部分支持**，见下)
* [x] 客户联系 (**大部分支持**，见下)
* [x] 微信客服 (**部分支持**，见下)
* [ ] 应用管理
* [x] 消息发送 (全部支持)
* [x] 消息接收
* [x] 素材管理 (**支持上传**, 见下)
* [x] OA (**大部分支持**，见下)
* [x] 会话内容存档 (**大部分支持**，见下)
* [x] 企业微信登录接口 (code2Session)
* [x] 获取访问用户身份 (code2UserInfo)

<details>
<summary>通讯录管理 API</summary>

* [ ] 成员管理
    - [ ] 创建成员
    - [x] 读取成员 *NOTE: 成员对外信息暂未实现*
    - [x] 更新成员
    - [ ] 删除成员
    - [ ] 批量删除成员
    - [ ] 获取部门成员
    - [x] 获取部门成员详情
    - [x] userid与openid互换
    - [ ] 二次验证
    - [ ] 邀请成员
    - [x] 获取加入企业二维码
    - [x] 手机号获取userid
    - [x] 邮箱获取userid
    - [ ] 获取成员ID列表
* [ ] 部门管理
    - [x] 创建部门
    - [ ] 更新部门
    - [ ] 删除部门
    - [x] 获取部门列表
    - [x] 获取子部门ID列表
* [ ] 标签管理
    - [ ] 创建标签
    - [ ] 更新标签名字
    - [ ] 删除标签
    - [ ] 获取标签成员
    - [ ] 增加标签成员
    - [ ] 删除标签成员
    - [ ] 获取标签列表
* [ ] 异步批量接口
    - [ ] 增量更新成员
    - [ ] 全量覆盖成员
    - [ ] 全量覆盖部门
    - [ ] 获取异步任务结果
* [ ] 通讯录回调通知
    - [ ] 成员变更通知
    - [ ] 部门变更通知
    - [ ] 标签变更通知
    - [ ] 异步任务完成通知

</details>

<details>
<summary>客户联系 API</summary>

* [x] 成员对外信息
* [x] 企业服务人员管理
    - [x] 获取配置了客户联系功能的成员列表
    - [x] 客户联系「联系我」管理
    - [x] 客户群「加入群聊」管理
* [x] 客户管理
    - [x] 获取客户列表
    - [x] 获取客户详情
    - [x] 批量获取客户详情
    - [x] 修改客户备注信息
* [x] 客户群管理
    - [x] 获取客户群列表
    - [x] 获取客户群详情
    - [x] 客户群opengid转换
* [x] 在职继承
    - [x] 分配在职成员的客户
    - [x] 查询客户接替状态
    - [ ] 分配在职成员的客户群
* [x] 离职继承
    - [ ] 获取待分配的离职成员列表
    - [x] 分配离职成员的客户
    - [x] 查询客户接替状态
    - [ ] 分配离职成员的客户群
* [x] 客户标签管理
   - [x] 管理企业标签
   - [x] 编辑客户企业标签
* [x] 客户分配
   - [x] 获取离职成员列表
   - [x] 分配在职或离职成员的客户
   - [x] 查询客户接替结果
   - [x] 分配离职成员的客户群
* [x] 变更回调通知
    - [x] 添加企业客户事件
    - [x] 编辑企业客户事件
    - [x] 外部联系人免验证添加成员事件
    - [x] 删除企业客户事件
    - [x] 删除跟进成员事件
    - [x] 客户接替失败事件
    - [x] 客户群变更事件
* [x] 消息推送
  - [x] 创建企业群发
  - [ ] 获取企业的全部群发记录
  - [x] 发送新客户欢迎语
  - [ ] 入群欢迎语素材管理

</details>

<details>
<summary>微信客服 API</summary>

* [x] 客服账号管理
    - [x] 添加客服账号
    - [x] 删除客服账号
    - [x] 修改客服账号
    - [x] 获取客服账号列表
    - [x] 获取客服账号链接
* [x] 接待人员管理
    - [x] 添加接待人员
    - [x] 删除接待人员
    - [x] 获取接待人员列表
* [x] 会话分配与消息收发
    - [x] 分配客服会话
    - [x] 接收消息和事件
    - [x] 发送消息
    - [x] 发送欢迎语等事件响应消息
* [ ] 「升级服务」配置
* [ ] 其他基础信息获取
    - [ ] 获取客户基础信息
* [ ] 统计管理
    - [ ] 获取「客户数据统计」企业汇总数据
    - [ ] 获取「客户数据统计」接待人员明细数据
* [ ] 机器人管理
    - [ ] 知识库分组管理
    - [ ] 知识库问答管理

</details>

<details>
<summary>身份验证 API</summary>

* [x] 获取访问用户身份

</details>

<details>
<summary>应用管理 API</summary>

* [ ] 获取应用
* [ ] 设置应用
* [ ] 自定义菜单
    - [ ] 创建菜单
    - [ ] 获取菜单
    - [ ] 删除菜单

</details>

<details>
<summary>消息发送 API</summary>

* [x] 发送应用消息
* [x] 接收消息
* [x] 发送消息到群聊会话
    - [x] 创建群聊会话
    - [x] 修改群聊会话
    - [x] 获取群聊会话
    - [x] 应用推送消息

### 消息类型

* [x] 文本消息
* [x] 图片消息
* [x] 语音消息
* [x] 视频消息
* [x] 文件消息
* [x] 文本卡片消息
* [x] 图文消息
* [x] 图文消息（mpnews）
* [x] markdown消息
* [x] 任务卡片消息

</details>

<details>
<summary>素材管理 API</summary>

* [x] 上传临时素材
* [x] 上传永久图片
* [ ] 获取临时素材
* [ ] 获取高清语音素材

</details>

<details>
<summary>OA API</summary>

* [ ] 审批
    - [x] 获取审批模板详情
    - [x] 提交审批申请
    - [x] 审批申请状态变化回调通知
    - [x] 批量获取审批单号
    - [x] 获取审批申请详情
    - [ ] 获取企业假期管理配置
    - [ ] 修改成员假期余额

</details>

<details>
<summary>会话内容存档 API</summary>

* [x] 获取会话内容存档开启成员列表
* [x] 获取会话同意情况
* [ ] 客户同意进行聊天内容存档事件回调
* [x] 获取会话内容存档内部群信息

</details>

## Notes

### 关于保密消息发送

Markdown 等类型消息目前不支持作为保密消息发送，强行发送会报错。
那么为何发送消息的方法还全部带着 `isSafe` 参数呢？

一方面，企业微信服务方完全可能在未来支持更多消息类型的保密发送，到时候不希望客户端代码重新编译；
另一方面，反正响应会报错，你也不会留着这种逻辑。因此不改了。

## License

* [MIT](./LICENSE)

## Stargazers over time

[![Stargazers over time](https://starchart.cc/xen0n/go-workwx.svg)](https://starchart.cc/xen0n/go-workwx)
