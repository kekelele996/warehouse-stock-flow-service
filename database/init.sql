-- ============================================================
-- wmsflow 供应链仓储出入库协同管理系统 - 数据库初始化脚本
-- 由 PostgreSQL 官方镜像在首次创建数据库时自动执行（幂等）
-- 说明：
--   1. 表结构由后端 GORM AutoMigrate 自动创建/同步，
--      完整 DDL 参考 backend/migrations/001_init.sql（手动初始化用）。
--   2. 演示数据（用户/货主/商品/库位/示例单据）由后端启动时 Seed 写入。
-- ============================================================

CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- 记录脚本已执行，便于排查
CREATE TABLE IF NOT EXISTS init_scripts (
    script_name VARCHAR(128) PRIMARY KEY,
    executed_at TIMESTAMPTZ DEFAULT now()
);
INSERT INTO init_scripts (script_name) VALUES ('init.sql')
ON CONFLICT (script_name) DO NOTHING;
