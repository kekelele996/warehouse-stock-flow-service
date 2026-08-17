# Bug 复现说明

## Bug 是什么

查看还没拣货的出库单、还没质检的入库单时，接口会触发空指针（nil pointer）崩溃。

## 如何触发

```bash
cd backend
go test ./...
```

## 错误信息

```
--- FAIL: TestInboundService_Receive_Success
    (panic) runtime error: invalid memory address or nil pointer dereference
--- FAIL: TestOutboundGet_PendingDoesNotPanic
    Get on pending outbound panicked: runtime error: invalid memory address or nil pointer dereference
--- FAIL: TestInboundGet_QCNotDoneDoesNotPanic
    Get on pending inbound panicked: runtime error: invalid memory address or nil pointer dereference
```
