# WMSFlow 供应链仓储出入库协同管理系统

面向供应链中转仓的多货主、多批次仓储出入库协同管理系统，覆盖收货质检、上架存储、拣货复核、出库发货全链路，强调货主与仓库之间的协同操作。后端 Go 1.22 + Gin + GORM，前端 Vue 3 + TypeScript + Vite，PostgreSQL 持久化，Redis 限流缓存。

## 快速启动（Docker Compose 一键部署）

```bash
docker compose up -d --build
```

启动完成后访问：

| 服务 | 地址 |
| --- | --- |
| 前端页面 | http://localhost:18709 |
| 后端健康检查 | http://localhost:19209/healthz |
| 后端 API | http://localhost:19209/api/v1 |
| PostgreSQL | localhost:44016（wmsflow_user / wmsflow_pwd） |
| Redis | localhost:46316 |

### 演示账号

| 账号 | 密码 | 角色 |
| --- | --- | --- |
| admin | wmsflow@123 | 系统管理员 |
| manager | wmsflow@123 | 仓库经理 |
| qc | wmsflow@123 | 质检员 |
| picker | wmsflow@123 | 拣货员 |
| checker | wmsflow@123 | 复核员 |
| owner01 | wmsflow@123 | 货主（华建建材） |

### 停止并清理

```bash
docker compose down -v --remove-orphans
```

## 项目主要功能

- **仓库总览**：今日收货/发货量、库位整体占用率环形图、各货主库存金额 TOP10、待处理任务列表。
- **入库管理**：入库单按状态 Tab 列表、创建入库单（选货主→添加商品明细）、收货质检（逐项录入实收数量与质检结果）、上架指派库位、状态流转 `Pending → Received → QCInProgress → Shelved → Completed`。
- **出库管理**：出库单列表、创建出库单（选货主→选商品→系统推荐库位）、拣货任务（逐项确认实拣数量并扣减库存）、复核、打包、发货、状态流转 `Pending → Picking → Checking → Packing → Shipped → Completed`。
- **库位管理**：按区域（A/B/C/D）网格视图（颜色表示占用率）、库位详情（当前存放商品）、批量创建库位。
- **货主管理**：货主档案、信用额度调整、暂停/恢复合作、货主详情（基本信息+库存统计+历史单据+账务概况）。
- **商品管理**：商品档案（SKU、条形码、品类、存储要求、单价等）。
- **操作日志**：收货/质检/上架/拣货/复核/发货等全链路操作审计。
- **RBAC 权限**：Admin / WarehouseManager / QCInspector / Picker / Checker / Owner，货主角色只能查看自己的单据与库存。

## 技术栈

| 层 | 技术栈 |
| --- | --- |
| 前端 | Vue 3 + TypeScript + Vite |
| 后端 | Go 1.22 + Gin + GORM |
| 数据库 | PostgreSQL 16 |
| 缓存 | Redis 7（限流，Redis 不可用时降级进程内限流） |
| 认证 | JWT + RBAC |
| 日志 | `log/slog` 结构化 JSON 日志 |
| 参数校验 | `github.com/go-playground/validator/v10` |
| 接口文档 | 见下文 API 清单（README） |

## 项目目录结构

```
wje-139/
├── backend/                    # Go 1.22 + Gin + GORM 后端
│   ├── cmd/server/main.go      # 启动入口：装配依赖、启动服务
│   ├── internal/
│   │   ├── config/             # 环境变量配置
│   │   ├── database/           # 数据库连接 + 自动迁移 + 种子数据
│   │   ├── model/              # 每个实体一个文件（user/owner/product/bin_location/inventory/inbound/outbound/operation_log）
│   │   ├── dto/                # 每个实体一个 DTO 文件（validator 校验 tag）
│   │   ├── repository/         # 每个实体一个仓储文件（GORM 实现 + 接口）
│   │   ├── service/            # 每个实体一个服务文件（业务状态机、事务）
│   │   ├── handler/            # 每个实体一个处理器文件
│   │   ├── router/             # 每个实体一个路由注册文件
│   │   ├── middleware/         # request_id/request_logger/error_handler/auth/rbac/audit/rate_limit
│   │   ├── constants/          # enums/error_codes/messages/log_templates/rbac/pagination
│   │   └── util/               # logger/jwt/password/response/app_error/formatters/order_no
│   ├── migrations/001_init.sql # 手动初始化完整 DDL
│   ├── pkg/pagination/         # 分页工具
│   ├── Dockerfile
│   └── go.mod / go.sum
├── frontend/                   # Vue 3 + TypeScript + Vite 前端
│   ├── src/
│   │   ├── api/                # 每个实体一个 API 文件
│   │   ├── components/         # StatusBadge/StepIndicator/StatCard/OccupancyRing/EmptyState/DataTable/ConfirmDialog/Modal
│   │   ├── pages/              # Login/Dashboard/Inbound/Outbound/BinLocations/Owners/OwnerDetail/Products/Audit/NotFound
│   │   ├── stores/             # 按实体拆分 Pinia store
│   │   ├── hooks/              # useAuth/usePagination/useInboundFlow/useOutboundFlow
│   │   ├── utils/              # request/format/generateOrderNo/toast
│   │   ├── constants/          # enums/roles
│   │   ├── types/              # 类型与枚举定义
│   │   ├── router/             # 路由与守卫
│   │   └── layouts/            # 主布局（侧边栏+顶栏）
│   ├── Dockerfile
│   └── nginx.conf              # SPA 路由 + /api 反向代理
├── database/init.sql           # 数据库初始化脚本
├── docker-compose.yml
├── .env
├── .env.example
└── README.md
```

## 本地开发

### 后端

```bash
cd backend
export GOPROXY=https://goproxy.cn,direct   # 网络受限时使用镜像
go mod tidy
go run ./cmd/server
```

构建命令：`go build ./...`；测试：`go test ./...`（仓储测试设置 `TEST_DATABASE_DSN` 连接真实库执行）。

### 前端

```bash
cd frontend
npm install --registry=https://registry.npmmirror.com
npm run dev        # http://localhost:5173（/api 代理到 http://localhost:19209）
npm run build      # 类型检查 + 产物构建
```

## 环境变量说明

| 变量 | 默认值 | 说明 |
| --- | --- | --- |
| COMPOSE_PROJECT_NAME | wmsflow | Compose 项目名，容器名带此前缀 |
| FRONTEND_PORT | 18709 | 前端宿主机端口（容器内 80） |
| BACKEND_PORT | 19209 | 后端宿主机端口（容器内 8080） |
| DB_PORT | 44016 | PostgreSQL 宿主机端口（容器内 5432） |
| REDIS_PORT | 46316 | Redis 宿主机端口（容器内 6379） |
| DB_NAME | wmsflow_db | 数据库名 |
| DB_USER | wmsflow_user | 数据库用户 |
| DB_PASSWORD | wmsflow_pwd | 数据库密码 |
| JWT_SECRET | change_me_to_a_long_random_string | JWT 密钥（≥16 位，生产必改） |
| JWT_EXPIRE_HOURS | 72 | JWT 有效期（小时） |
| LOG_LEVEL | info | 日志级别 |

## API 清单

统一前缀 `/api/v1`，健康检查 `/healthz`，响应格式 `{ "code": 0, "message": "ok", "data": ... }`，分页参数 `page` / `page_size`。

### 认证

| 方法 | 路径 | 说明 |
| --- | --- | --- |
| POST | /api/v1/auth/register | 注册货主账号（同步创建货主档案） |
| POST | /api/v1/auth/login | 登录，返回 JWT |
| GET | /api/v1/auth/me | 当前用户信息 |

### 货主 Owner

| 方法 | 路径 | 说明 |
| --- | --- | --- |
| GET | /api/v1/owners | 货主列表（分页/关键词/状态） |
| POST | /api/v1/owners | 创建货主（owner:manage） |
| GET | /api/v1/owners/:id | 货主详情 |
| GET | /api/v1/owners/:id/stats | 货主统计（商品/单据/库存金额） |
| PUT | /api/v1/owners/:id | 更新货主 |
| PUT | /api/v1/owners/:id/credit | 调整信用额度 |
| PUT | /api/v1/owners/:id/suspend | 暂停合作 |
| PUT | /api/v1/owners/:id/activate | 恢复合作 |

### 商品 Product

| 方法 | 路径 | 说明 |
| --- | --- | --- |
| GET | /api/v1/products | 商品列表（货主/关键词/品类/存储要求过滤） |
| POST | /api/v1/products | 创建商品 |
| GET | /api/v1/products/by-owner/:ownerId | 按货主查询商品（复用商品列表逻辑） |
| GET | /api/v1/products/:id | 商品详情 |
| PUT | /api/v1/products/:id | 更新商品 |

### 库位 BinLocation

| 方法 | 路径 | 说明 |
| --- | --- | --- |
| GET | /api/v1/bin-locations | 库位列表（区域/状态过滤） |
| POST | /api/v1/bin-locations | 创建库位（bin:manage） |
| POST | /api/v1/bin-locations/batch | 批量创建库位（bin:manage） |
| POST | /api/v1/bin-locations/recommend | 推荐库位（bin:manage） |
| GET | /api/v1/bin-locations/:id | 库位详情 |
| GET | /api/v1/bin-locations/:id/contents | 库位存放商品 |
| PUT | /api/v1/bin-locations/:id | 更新库位（bin:manage） |

### 库存 Inventory

| 方法 | 路径 | 说明 |
| --- | --- | --- |
| GET | /api/v1/inventory | 库存列表（货主/商品/库位过滤） |
| GET | /api/v1/inventory/summary | 货主库存金额汇总 |

### 入库 Inbound

| 方法 | 路径 | 说明 |
| --- | --- | --- |
| GET | /api/v1/inbound | 入库单列表（状态/货主过滤） |
| POST | /api/v1/inbound | 创建入库单（含明细） |
| GET | /api/v1/inbound/:id | 入库单详情（含明细） |
| POST | /api/v1/inbound/:id/receive | 收货 Pending→Received（Admin/WarehouseManager） |
| POST | /api/v1/inbound/:id/qc | 质检 Received→QCInProgress（+QCInspector） |
| POST | /api/v1/inbound/:id/shelve | 上架 QCInProgress→Shelved（写库存） |
| POST | /api/v1/inbound/:id/complete | 完成 Shelved→Completed |

### 出库 Outbound

| 方法 | 路径 | 说明 |
| --- | --- | --- |
| GET | /api/v1/outbound | 出库单列表（状态/货主过滤） |
| POST | /api/v1/outbound | 创建出库单（含明细） |
| GET | /api/v1/outbound/:id | 出库单详情（含明细） |
| POST | /api/v1/outbound/:id/picking | 拣货 Pending→Picking（扣减库存，+Picker） |
| POST | /api/v1/outbound/:id/checking | 复核 Picking→Checking（+Checker） |
| POST | /api/v1/outbound/:id/packing | 打包 Checking→Packing |
| POST | /api/v1/outbound/:id/ship | 发货 Packing→Shipped（快递单号） |
| POST | /api/v1/outbound/:id/complete | 完成 Shipped→Completed |

### 总览与审计

| 方法 | 路径 | 说明 |
| --- | --- | --- |
| GET | /api/v1/dashboard/summary | 仓库总览汇总 |
| GET | /api/v1/audit | 操作日志列表（audit:view） |

## API 调用示例（curl）

```bash
# 1. 登录获取 JWT
TOKEN=$(curl -s http://localhost:19209/api/v1/auth/login \
  -H 'Content-Type: application/json' \
  -d '{"username":"admin","password":"wmsflow@123"}' | python3 -c 'import sys,json;print(json.load(sys.stdin)["data"]["token"])')

# 2. 健康检查
curl -s http://localhost:19209/healthz

# 3. 货主列表（带 JWT）
curl -s http://localhost:19209/api/v1/owners?page=1\&page_size=5 \
  -H "Authorization: Bearer $TOKEN"

# 4. 创建货主
curl -s http://localhost:19209/api/v1/owners \
  -H "Authorization: Bearer $TOKEN" -H 'Content-Type: application/json' \
  -d '{"name":"新货主有限公司","contact_name":"王五","phone":"13900001111","email":"wangwu@example.com","settlement_method":"Monthly","credit_limit":100000}'

# 5. 创建入库单
curl -s http://localhost:19209/api/v1/inbound \
  -H "Authorization: Bearer $TOKEN" -H 'Content-Type: application/json' \
  -d '{"owner_id":1,"supplier_name":"示例供应商","items":[{"product_id":1,"batch_no":"B20260801","expected_qty":100}]}'

# 6. 收货（流转 Pending->Received）
curl -s -X POST http://localhost:19209/api/v1/inbound/1/receive \
  -H "Authorization: Bearer $TOKEN"

# 7. 仓库总览
curl -s http://localhost:19209/api/v1/dashboard/summary -H "Authorization: Bearer $TOKEN"
```

## Docker 部署说明

- **端口映射**：前端 `${FRONTEND_PORT:-18709}:80`，后端 `${BACKEND_PORT:-19209}:8080`，数据库 `${DB_PORT:-44016}:5432`，Redis `${REDIS_PORT:-46316}:6379`。
- **数据卷**：`db_data` 持久化 PostgreSQL 数据，`redis_data` 持久化 Redis 数据（命名卷，兼容任意目录名）。
- **健康检查**：db 使用 `pg_isready`，backend 使用 `wget /healthz`，backend 通过 `depends_on: condition: service_healthy` 等待数据库就绪。
- **中文目录名**：Compose 使用命名卷与环境变量注入，无绑定挂载中文路径，可在任意目录名下启动。
- **常见问题**：
  1. 端口冲突：修改 `.env` 中的 `*_PORT` 后重新 `docker compose up -d`。
  2. 数据库数据重置：`docker compose down -v` 会删除命名卷数据，如需保留请勿加 `-v`。
  3. JWT 密钥：生产环境务必修改 `.env` 中 `JWT_SECRET`（≥16 位）。

## 枚举出现位置清单

| 枚举 | 后端出现位置 | 前端出现位置 |
| --- | --- | --- |
| 入库单状态 InboundStatus（Pending/Received/QCInProgress/Shelved/Completed） | `backend/internal/constants/enums.go`、`model/inbound_order.go`、`dto/inbound_dto.go`、`service/inbound_service.go`（状态机）、`handler/inbound_handler.go`、`repository/inbound_repository.go`、`util/formatters.go`、`constants/log_templates.go`、`constants/messages.go`、`middleware/rbac.go` | `frontend/src/types/enums.ts`、`constants/enums.ts`、`hooks/useInboundFlow.ts`、`api/inbound.ts`、`pages/InboundPage.vue`、`stores/inbound.ts` |
| 出库单状态 OutboundStatus（Pending/Picking/Checking/Packing/Shipped/Completed） | `backend/internal/constants/enums.go`、`model/outbound_order.go`、`dto/outbound_dto.go`、`service/outbound_service.go`（状态机）、`handler/outbound_handler.go`、`repository/outbound_repository.go`、`util/formatters.go`、`constants/log_templates.go`、`constants/messages.go` | `frontend/src/types/enums.ts`、`constants/enums.ts`、`hooks/useInboundFlow.ts`、`api/outbound.ts`、`pages/OutboundPage.vue`、`stores/outbound.ts` |
| 存储要求 StorageRequirement（Normal/ColdChain/Dangerous/Fragile） | `backend/internal/constants/enums.go`、`model/product.go`、`model/bin_location.go`、`dto/product_dto.go`、`dto/bin_location_dto.go`、`service/bin_location_service.go`（推荐库位）、`repository/bin_location_repository.go`、`util/formatters.go`、`constants/log_templates.go` | `frontend/src/types/enums.ts`、`constants/enums.ts`、`api/product.ts`、`api/binLocation.ts`、`pages/BinLocationsPage.vue`、`pages/ProductsPage.vue` |
| 质检结果 QCResult（Pass/Fail/Partial） | `backend/internal/constants/enums.go`、`model/inbound_order.go`、`dto/inbound_dto.go`、`service/inbound_service.go`、`util/formatters.go` | `frontend/src/types/enums.ts`、`constants/enums.ts`、`api/inbound.ts`、`pages/InboundPage.vue` |
| 结算方式 SettlementMethod（Monthly/PerOrder/Prepaid） | `backend/internal/constants/enums.go`、`model/owner.go`、`dto/owner_dto.go`、`util/formatters.go` | `frontend/src/types/enums.ts`、`constants/enums.ts`、`pages/OwnersPage.vue` |
| 货主状态 OwnerStatus（Active/Suspended） | `backend/internal/constants/enums.go`、`model/owner.go`、`service/owner_service.go`、`service/inbound_service.go`、`service/outbound_service.go`、`util/formatters.go` | `frontend/src/types/enums.ts`、`constants/enums.ts`、`pages/OwnersPage.vue` |
| 库位状态 BinLocationStatus（Available/Occupied/Reserved/Maintenance） | `backend/internal/constants/enums.go`、`model/bin_location.go`、`service/inventory_service.go`（重算占用）、`util/formatters.go` | `frontend/src/types/enums.ts`、`constants/enums.ts`、`pages/BinLocationsPage.vue` |
| 角色 UserRole（Admin/WarehouseManager/QCInspector/Picker/Checker/Owner） | `backend/internal/constants/rbac.go`、`model/user.go`、`middleware/rbac.go`、`middleware/auth.go`、`util/jwt.go`、`dto/auth_dto.go`、`util/formatters.go`、`constants/log_templates.go`、`constants/messages.go`、`database/seed.go` | `frontend/src/types/enums.ts`、`constants/enums.ts`、`constants/roles.ts`、`hooks/useAuth.ts`、`router/index.ts`、`layouts/MainLayout.vue`、`pages/*` |

## 说明

- 商品新增 `price`（单价）字段用于支撑「各货主库存金额」统计（提示词原字段清单未列单价，属合理扩展）。
- 库存金额 = Σ(库存数量 × 商品单价)，库位占用率 = Σ(商品体积 × 数量) / 库位容量。
- 并发安全：拣货扣减库存使用 `SELECT ... FOR UPDATE` 行锁；入库上架、出库扣减均置于 service 层事务。
- 本仓库严格遵守「严禁合并职责到单一文件」：每个实体的 model/dto/repository/service/handler/router/constants 均独立文件。

## License

MIT
