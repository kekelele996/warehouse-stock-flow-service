# Bug 复现说明

## Bug 是什么

入库质检后上架时，质检不合格的商品可以上架、合格商品反而被拦截；上架写进库存的数量用了「期望数量」而不是「实收数量」，明细视图里的实收数量与质检结果文案也错误；同时入库单状态机允许从 `QCInProgress` 直接跳到 `Completed`。

## 如何触发

```bash
cd backend
go test -run 'TestInboundShelve_PassAddsActualQty|TestInboundShelve_FailRejected|TestInboundComplete_FromQCRejected' ./... -count=1
```

## 错误信息

```
--- FAIL: TestInboundShelve_PassAddsActualQty
    shelve error: code=40900 message=质检不合格的商品不能上架（明细ID 11）
--- FAIL: TestInboundShelve_FailRejected
    expected fail item to be rejected
--- FAIL: TestInboundComplete_FromQCRejected
    expected transition error
```
