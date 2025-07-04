
---

# ✅ Golang 配置服务设计文档（支持灰度、版本控制、策略路由、模块组）

---

## 一、📌 项目背景

我当前负责的 PC 客户端 SDK 主要实现配置下发功能，配置流程分为以下三个阶段：

1. **嗅探接口**：客户端上传平台/版本/配置 Hash，服务端判断是否有新版本，若有返回 CDN 地址；
2. **静态配置 CDN 下载**：支持全量和增量配置拉取，配置结构为模块组下属多个模块；
3. **动态配置接口**：客户端上传自身特征（平台、版本、渠道等），服务端返回模块组 → 模块 ID 映射表；
4. **客户端本地组装**：通过模块 ID 在静态配置中提取内容，最终生成模块组级别配置。

该架构可支持灰度发布、AB 测试、策略化分发，模块生命周期管理，以及静态配置的高性能缓存。

---

## 二、🧠 核心设计理念

| 项目     | 说明                            |
| ------ | ----------------------------- |
| 配置单元粒度 | 以模块为最小单元，按模块组聚合               |
| 策略灵活性  | 动态规则根据客户端特征匹配模块 ID，支持优先级与灰度路由 |
| 端侧负担   | 配置筛选和组装逻辑保持在客户端，服务端返回决策结果     |
| 更新机制   | 服务端增量打包静态配置，减少带宽；版本由客户端嗅探决定   |
| 生命周期管理 | 每个模块支持启用状态和有效期，便于限时配置或清理      |

---

## 三、🧩 数据模型设计

### 1. 模块组（config\_module\_groups）

| 字段                        | 类型        | 说明       |
| ------------------------- | --------- | -------- |
| group\_id                 | string    | 模块组唯一 ID |
| name                      | string    | 模块组名称    |
| description               | text      | 说明       |
| created\_at / updated\_at | timestamp | 时间戳      |

### 2. 模块（config\_modules）

| 字段                        | 类型        | 说明      |
| ------------------------- | --------- | ------- |
| module\_id                | string    | 模块唯一 ID |
| group\_id                 | string    | 所属模块组   |
| name                      | string    | 模块名称    |
| content                   | JSON      | 配置内容    |
| valid\_from / valid\_to   | datetime  | 有效期     |
| enabled                   | bool      | 是否启用    |
| version                   | int       | 模块版本号   |
| created\_at / updated\_at | timestamp | 时间戳     |

### 3. 动态策略规则（config\_rules）

| 字段           | 类型     | 说明                                |
| ------------ | ------ | --------------------------------- |
| rule\_id     | string | 规则 ID                             |
| group\_id    | string | 所属模块组                             |
| module\_id   | string | 命中模块 ID                           |
| filter\_expr | string | 条件表达式（如 channel=A & version>=2.1） |
| priority     | int    | 优先级                               |
| enabled      | bool   | 是否启用                              |

---

## 四、📡 核心接口说明

### 1. 嗅探接口（check\_update）

```http
GET /api/configs/check_update?platform=win&version=2.1.0&hash=abc123

返回：
{
  "need_update": true,
  "cdn_url": "https://cdn.xxx.com/conf/v3.json",
  "update_mode": "diff",
  "new_hash": "def456"
}
```

> ✅ 用于判断是否需要重新拉取静态配置文件。

---

### 2. 静态配置 CDN 内容结构（例）

```json
{
  "module_abc": { "name": "实验版首页", "content": { ... } },
  "module_def": { "name": "默认首页", "content": { ... } },
  "module_dark": { "name": "深色模式", "content": { ... } }
}
```

> ✅ 客户端拉取后会缓存此 JSON，在动态请求后按模块 ID 提取配置。

---

### 3. 动态配置接口（filter\_ids）

```http
POST /api/configs/filter_ids
{
  "platform": "win",
  "version": "2.1.0",
  "channel": "A"
}

返回：
{
  "homepage": "module_abc",
  "search": "module_default",
  "sidebar": "module_exp_2"
}
```

> ✅ 每个模块组只命中一个模块 ID，客户端本地组装配置。

---

### 4. 客户端组装示意（伪代码）

```cpp
auto dynamic_ids = get_dynamic_ids(); // 从接口获取模块组到模块ID的映射
auto static_config = load_cdn_config(); // 加载静态配置大 JSON

map<string, JSON> result;
for (auto &[group, module_id] : dynamic_ids) {
    result[group] = static_config[module_id];
}
```

---

## 五、🧪 A/B 测试策略示例

| group\_id | 条件表达式         | 命中模块            |
| --------- | ------------- | --------------- |
| homepage  | 渠道=A 且版本>=2.0 | module\_abc     |
| homepage  | 渠道=B          | module\_b\_test |
| homepage  | default       | module\_default |

在策略路由生效后：

* 客户端 A（渠道 A）命中实验模块；
* 客户端 B（渠道 B）命中 B 测试模块；
* 其他命中默认模块。

---

