# Bug 复现说明

## Bug 是什么

出库单拣货时，库存扣减数量使用了「期望数量」而不是「实际拣货数量」，拣货完成后明细里展示的实拣数量也错了；同时出库单状态机允许从 `Pending` 直接跳到 `Checking`，且拣货完成后的状态文案显示错误。

## 如何触发

```bash
cd backend
go test ./internal/service/ -run 'TestOutboundPicking_Composite|TestOutboundChecking_FromPendingRejected' -count=1
```

## 错误信息

```
--- FAIL: TestOutboundPicking_Composite
    expected reduce once with actual qty 3, got calls=1 qty=5
--- FAIL: TestOutboundChecking_FromPendingRejected
    expected transition error
```
