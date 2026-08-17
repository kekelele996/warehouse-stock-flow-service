# Bug 复现说明

## Bug 是什么

货主角色登录后，库存列表、入库单列表、出库单列表都没有按当前货主做数据范围过滤，能看到其他货主的数据。

## 如何触发

```bash
cd backend
go test ./...
```

## 错误信息

```
--- FAIL: TestInventoryService_List_OwnerScope
    expected scoped owner 6, got 0
--- FAIL: TestInventoryList_OwnerScopeApplied
    expected scoped owner 6, got 0
--- FAIL: TestInboundList_OwnerScopeApplied
    expected scoped owner 6, got 0
--- FAIL: TestOutboundList_OwnerScopeApplied
    expected scoped owner 6, got 0
```
