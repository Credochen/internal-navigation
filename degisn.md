针对团队内部多系统（多为 IP 地址形式）的导航需求，设计一个轻量、高效且易于维护的网址导航网站。该方案完全基于 **Go + SQLite + Vue 3** 架构，主打“开箱即用”与“极简运维”。

---

## 1. 系统架构与目录结构

为了保证内部工具的轻量化，采用**单兵作战/极简部署**模式。后端将 SQLite 嵌入在服务中，前端打包后的静态资源直接通过 Go 的 `embed` 特性嵌入到二进制文件中。团队部署时，**只需运行一个二进制文件即可启动完整服务**。

### 项目目录结构

```text
internal-navigation/
├── backend/
│   ├── cmd/
│   │   └── main.go          # 程序入口
│   ├── config/              # 配置管理
│   ├── handler/             # 路由处理器 (API)
│   ├── model/               # SQLite 结构体定义与初始化
│   └── repository/          # 数据库增删改查
├── frontend/
│   ├── src/
│   │   ├── assets/          # 静态资源（图标等）
│   │   ├── components/      # 抽离的组件（卡片、弹窗）
│   │   └── App.vue          # 主页面
│   ├── package.json
│   └── vite.config.js
├── dist/                    # 前端打包输出（Go embed 静态指向此处）
└── go.mod

```

---

## 2. 数据库设计 (SQLite)

针对团队内部场景，设计两张核心表：**分类表**（如：基础架构、数据工具、CI/CD）和**网站条目表**。

### `categories` (分类表)

```sql
CREATE TABLE IF NOT EXISTS categories (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE,          -- 分类名称，如 "CI/CD 平台"
    sort_order INTEGER DEFAULT 0,       -- 排序权重，越小越靠前
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

```

### `navigations` (导航条目表)

```sql
CREATE TABLE IF NOT EXISTS navigations (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    category_id INTEGER NOT NULL,       -- 关联分类
    title TEXT NOT NULL,                -- 系统名称，如 "KafkaMap"
    url TEXT NOT NULL,                  -- 访问地址，如 "http://10.0.0.15:8080"
    icon TEXT,                          -- 图标（支持放 SVG、开源图标库类名或首字母）
    description TEXT,                   -- 备注/用途说明，如 "Kafka 消息队列监控"
    sort_order INTEGER DEFAULT 0,       -- 排序权重
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY(category_id) REFERENCES categories(id) ON DELETE CASCADE
);

```

---

## 3. 后端设计 (Go)

后端选用轻量级 Web 框架（如 Gin 或 Fiber）。鉴于团队内部使用，初期无需复杂的权限系统，可设计一个简单的 `Admin` 路由组，通过配置文件中的固定 Token 或简易密码进行拦截，用以管理（增删改查）导航数据。

### 核心 API 接口设计

| 请求方法 | 路由 | 功能说明 | 权限 |
| --- | --- | --- | --- |
| **GET** | `/api/v1/navs` | 获取完整的分类及导航数据（树状结构） | 公开 |
| **POST** | `/api/v1/admin/categories` | 创建新分类 | 管理员 |
| **PUT/DELETE** | `/api/v1/admin/categories/:id` | 修改/删除分类 | 管理员 |
| **POST** | `/api/v1/admin/navs` | 创建导航条目 | 管理员 |
| **PUT/DELETE** | `/api/v1/admin/navs/:id` | 修改/删除导航条目 | 管理员 |

### 核心数据响应格式 (Tree Structure)

`GET /api/v1/navs` 返回的数据直接按分类聚合，方便前端一键渲染：

```json
[
  {
    "category_id": 1,
    "category_name": "数据与中间件",
    "items": [
      {
        "id": 10,
        "title": "KafkaMap",
        "url": "http://10.20.1.5:9001",
        "icon": "polyline",
        "description": "测试环境 Kafka 集群流向监控"
      }
    ]
  }
]

```

---

## 4. 前端设计 (Vue 3)

### 界面视觉与交互风格

* **极简大气布局**：采用**左侧/顶部固定的分类导航栏 + 右侧/下方平铺的卡片流（Card Grid）**。支持点击分类锚点平滑滚动。
* **搜索框高亮**：顶部常驻一个支持快捷键（如 `Ctrl + K`）的搜索框。输入系统名称、URL 甚至描述的关键字时，瞬间高亮并过滤相关卡片。对记不住 IP 的团队成员极其友好。
* **暗黑模式支持**：团队研发多为“夜猫子”，原生支持浅色/深色模式一键切换。
* **卡片视觉设计**：
* **Hover 动效**：鼠标悬停在卡片上时，卡片轻微上浮并出现阴影变色，强调可点击性。
* **标签化展示**：在卡片右下角或右上角通过小标签（Badge）展示其环境属性（如 `Prod`、`Test`、`Dev`）。



### 前端核心组件技术选型

* **基础框架**：Vue 3 (Setup 语法糖) + Vite
* **UI 组件库**：Tailwind CSS（极易定制出简洁、高级感的纯净界面，避免传统组件库的厚重感）+ Element Plus 或 Ant Design Vue（用于后台管理弹窗）。
* **图标库**：Lucide Vue 或 FontAwesome（提供技术栈相关的极简线条图标）。

---

## 5. 提效与亮点功能（针对内部场景优化）

1. **内网探活 (Health Check)**
* **功能**：由 Go 后端在后台起一个轻量级定时任务（如每 5 分钟），对所有配置的 IP 地址进行 `Ping` 或 `HTTP HEAD` 探测。
* **前端呈现**：在导航卡片的标题旁，显示一个小绿点（在线）或小红点（离线）。**系统挂没挂，看一眼导航就知道**。


2. **静态资源内嵌 (`go:embed`)**
* 前端 `npm run build` 生成的静态文件，通过 Go 1.16+ 的 `//go:embed dist/*` 直接打入 Go 二进制中。
* 运维时不需要配置 Nginx 转发静态资源，直接 `nohup ./internav &` 即可全栈跑通。


3. **一键快速导入导出**
* 支持通过 JSON 或 Excel 批量导入团队现有的 IP 资产清单，实现一分钟快速建站。