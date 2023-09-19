<img align="right" src="./SQLE_logo.png">

简体中文 | [English](./README_en.md)

[![Release](https://img.shields.io/github/release/actiontech/sqle.svg?style=flat-square)](https://github.com/actiontech/sqle/releases)
[![GitHub license](https://img.shields.io/github/license/actiontech/sqle.svg)](https://github.com/actiontech/sqle/blob/main/LICENSE)
[![GitHub stars](https://img.shields.io/github/stars/actiontech/sqle.svg)](https://github.com/actiontech/sqle/stargazers)
[![GitHub issues](https://img.shields.io/github/issues/actiontech/sqle.svg)](https://github.com/actiontech/sqle/issues)
[![GitHub closed issues](https://img.shields.io/github/issues-closed-raw/actiontech/sqle.svg)](https://github.com/actiontech/sqle/issues?q=is%3Aissue+is%3Aclosed)
[![Docker Pulls](https://img.shields.io/docker/pulls/actiontech/sqle-ce.svg)](https://hub.docker.com/r/actiontech/sqle-ce)

SQLE 由上海爱可生信息技术股份有限公司（以下简称爱可生公司）出品和维护，是爱可生公司[云树®SQL质量管理软件SQLE](https://www.actionsky.com/sqle)（简称：CTREE SQLE）软件产品的开源版本。SQLE 是一个支持多场景，原生支持 MySQL 审核且数据库类型可扩展的 SQL 审核工具。

[官方网站](https://opensource.actionsky.com/sqle/) | [文档](https://actiontech.github.io/sqle-docs/docs/intro/) | [安装](https://actiontech.github.io/sqle-docs/docs/deploy-manual/intro) | [在线体验](https://actiontech.github.io/sqle-docs/docs/online-demo)

## 产品展示
![product_show](./SQLE_product_show.gif)

[更多产品展示](https://actiontech.github.io/sqle-docs-cn/0.overview/2_product_show.html)

## 产品特色

|特色|说明|
|--|--|
|SQL审核规范 |1. 审核规则自定义（700+）<br/>2. 支持审核结果分级展示，支持生成下载审核报告<br/>3. 支持规则模版，灵活组合规则<br/>4. 审核白名单，跳过特例 SQL<br/>5. 支持集成 IDE 自助审核<br/> |
| 多场景审核 | 支持事前事后审核，覆盖开发、测试、上线、生产等环节
| 标准化上线流程 | 1. SQL 审核流程按需自定义，满足企业内部不同流程管理要求<br/>2. 支持定时上线<br/>3. 支持设置运维时间<br/>4. 支持 Online DDL |
| 多数据库类型支持 | 1. 统一接口，可通过插件进行多数据库审核扩展<br/> 2. 内置 MySQL 审核插件，官方支持常用数据库类型（包括：PostgreSQL、DB2、SQL Server、Oracle、TiDB、OceanBase） |
| 统一的 SQL 客户端入口 | 提供审核管控的 SQL 客户端，杜绝执行不合规 SQL|
| 丰富的集成能力 | 1. 标准 HTTP API 接口可与客户内部流程系统对接<br/> 2. 支持 LDAP，OAuth2.0 用户对接<br/> 3. 支持邮件、微信企业号、webhook 告警对接 |

## 应用场景
|场景|介绍|
| --- | --- |
| 上线前控制 | 1. SQLE 赋能开发，在代码开发阶段检查 SQL 质量<br/> 2. SQLE 集成专家经验，形成可复用的 SQL 规范标准规则，用户在平台提交工单后，平台将基于审核规则模板对提交的 SQL 语句进行初审，用以解决事前审核规范不标准难题|
|上线后监督| 1. SQLE 提供智能扫描审核功能，实现生产环境下的 SQL 审核优化 <br/> 2. SQLE 支持多种类型的扫描任务，基于任务需求执行周期性的扫描任务，并生成扫描结果报告，及时告警|

## 在线体验

| 社区版 | 企业版 |
| --- | --- |
| [SQLE 社区版](http://demo.sqle.actionsky.com/) | [SQLE 企业版](http://demo.sqle.actionsky.com:8889/) |
| 超级管理员: admin <br/> 密码： admin | 超级管理员: admin <br/> 密码： admin | 

### 测试 MySQL
|配置项|值|
|---|---|
| 地址 | 20.20.20.3 |
| 端口 | 3306 |
| 用户 | root |
| 密码 | test 

> 注意事项
> 1. 该服务仅用于在线功能体验，请勿在生产环境使用；
> 2. 该测试服务数据会定期清理。

## SQL 审核插件
目前支持其他种类数据库的审核插件:
* [PostgreSQL](https://github.com/actiontech/sqle-pg-plugin)
* [Oracle](https://github.com/actiontech/sqle-oracle-plugin)
* [SQL Server](https://github.com/actiontech/sqle-ms-plugin)
* [DB2](https://github.com/actiontech/sqle-db2-plugin)

更多了解：《[功能说明及开发手册](https://actiontech.github.io/sqle-docs/docs/dev-manual/plugins/intro) 》

## 官方技术支持

|渠道 | 链接 |
| -- | -- |
| 代码库 | [github.com/actiontech/sqle](https://github.com/actiontech/sqle) |
| UI 库 | [github.com/actiontech/sqle-ui](https://github.com/actiontech/sqle-ui) |
| 文档库 | [github.com/actiontech/sqle-docs](https://github.com/actiontech/sqle-docs) |
| 文档主页 | [actiontech.github.io/sqle-docs](https://actiontech.github.io/sqle-docs/) |
| 社区网站 | [opensource.actionsky.com](https://opensource.actionsky.com) |
| 微信技术交流群 | 添加管理员：ActionOpenSource |
| 开源社区微信公众号 | ![QR_code](./QR_code.png) |

## 联系我们
如果想获得 SQLE 的商业支持, 您可以联系我们：
* 全国支持: 400-820-6580
* 华北地区: 86-13910506562, 汪先生
* 华南地区: 86-18503063188, 曹先生
* 华东地区: 86-18930110869, 梁先生
* 西南地区: 86-13540040119, 洪先生
