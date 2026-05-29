# Internal Navigation

轻量级团队内部网址导航系统，基于 **Go + SQLite + Vue 3** 架构，主打“开箱即用”与“极简运维”。

部署时**只需运行一个二进制文件**即可启动完整服务，前端静态资源已通过 `go:embed` 嵌入二进制中。

---

## 特性

- **单文件部署**：Go 二进制内嵌前端静态资源，无需 Nginx
- **树状导航**：按分类聚合展示，支持分类锚点平滑滚动
- **实时搜索**：支持按系统名称、URL、描述搜索（快捷键 `Ctrl + K`）
- **健康探活**：后台定时探测各系统在线状态，卡片显示绿点/红点
- **暗黑模式**：一键切换浅色/深色主题
- **管理后台**：弹窗式分类与导航条目管理，支持 JSON 导入导出
- **轻量鉴权**：通过 `X-Admin-Token` Header 进行简易管理员鉴权

---

## 技术栈

| 层级 | 技术 |
|------|------|
| 后端 | Go 1.22 + Gin + GORM + SQLite |
| 前端 | Vue 3 + Vite + Tailwind CSS |
| 部署 | `go:embed` 静态资源嵌入 |

---

## 快速开始

### 1. 克隆仓库

```bash
git clone <repo-url>
cd internal-navigation
```

### 2. 构建（本地）

```bash
make build
```

### 3. 运行

```bash
./internal-navigation
```

默认监听 `http://localhost:8080`。

---

## 交叉编译 Linux 二进制

在 macOS 或 Windows 开发环境下，编译可在 Linux 服务器运行的二进制：

### 方式一：禁用 CGO（推荐，无需额外依赖）

```bash
make build-linux-simple
```

生成的 `internal-navigation-linux-amd64` 可直接在 Linux x86_64 服务器运行。

### 方式二：启用 CGO + 静态链接（需要 musl-cross 工具链）

```bash
# macOS 安装交叉编译工具链
brew install FiloSottile/musl-cross/musl-cross

make build-linux
```

---

## systemd 部署（Linux）

### 1. 准备目录与二进制

```bash
sudo mkdir -p /var/lib/internal-navigation
sudo cp internal-navigation-linux-amd64 /usr/local/bin/internal-navigation
sudo chmod +x /usr/local/bin/internal-navigation
```

### 2. 修改默认 Token

```bash
sudo sed -i 's/changeme/YOUR_SECURE_TOKEN/g' /etc/systemd/system/internal-navigation.service
```

> 或直接编辑 `internal-navigation.service` 文件中的 `ADMIN_TOKEN`。

### 3. 安装 systemd 服务

```bash
sudo cp internal-navigation.service /etc/systemd/system/
sudo systemctl daemon-reload
sudo systemctl enable internal-navigation
sudo systemctl start internal-navigation
```

### 4. 查看状态

```bash
sudo systemctl status internal-navigation
sudo journalctl -u internal-navigation -f
```

---

## 配置说明

通过环境变量配置：

| 变量 | 默认值 | 说明 |
|------|--------|------|
| `PORT` | `8080` | 服务监听端口 |
| `ADMIN_TOKEN` | `admin123` | 管理后台鉴权 Token |
| `DB_PATH` | `nav.db` | SQLite 数据库文件路径 |

示例：

```bash
PORT=80 ADMIN_TOKEN=mytoken ./internal-navigation
```

---

## API 概览

### 公开接口

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/navs` | 获取分类导航树（含在线状态） |
| GET | `/api/v1/categories` | 获取分类列表 |

### 管理接口（需 `X-Admin-Token`）

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/v1/admin/categories` | 创建分类 |
| PUT | `/api/v1/admin/categories/:id` | 更新分类 |
| DELETE | `/api/v1/admin/categories/:id` | 删除分类 |
| POST | `/api/v1/admin/navs` | 创建导航 |
| PUT | `/api/v1/admin/navs/:id` | 更新导航 |
| DELETE | `/api/v1/admin/navs/:id` | 删除导航 |
| GET | `/api/v1/admin/export` | 导出 JSON |
| POST | `/api/v1/admin/import` | 导入 JSON |

---

## 导入导出数据格式

```json
{
  "categories": [
    { "id": 1, "name": "基础架构", "sort_order": 0 }
  ],
  "navigations": [
    { "id": 1, "category_id": 1, "title": "Prometheus", "url": "http://10.0.0.20:9090", "icon": "📈", "description": "监控系统", "sort_order": 0 }
  ]
}
```

---

## 目录结构

```
internal-navigation/
├── main.go                       # 程序入口
├── Makefile                      # 构建脚本
├── internal-navigation.service   # systemd 服务文件
├── backend/
│   ├── config/      # 配置管理
│   ├── handler/     # HTTP 处理器
│   ├── model/       # 数据模型
│   └── repository/  # 数据访问层
├── frontend/
│   ├── src/         # Vue 3 源码
│   ├── index.html
│   ├── package.json
│   └── vite.config.js
└── dist/            # 前端构建产物（被 Go embed 嵌入）
```

---

## License

MIT
