# Bug 复现说明

## Bug 是什么

给商品推荐库位时选错了库位（选了占用率最高的而不是最低的）；返回的库位编码、状态文字、存储要求文字都不正确。

## 如何触发

```bash
cd backend
go test ./... -run 'TestBinRecommend_PicksLowestOccupancy|TestBinView_CodeAndStatusText' -count=1
```

## 错误信息

```
--- FAIL: TestBinRecommend_PicksLowestOccupancy
    expected lowest occupancy bin id 1, got 2 (occupancy 80)
--- FAIL: TestBinView_CodeAndStatusText
    unexpected bin code ...
    unexpected status text ...
```
