// 单号生成工具：INB-YYYYMMDD-XXXX / OUT-YYYYMMDD-XXXX（与后端单号规则一致）。
export function generateOrderNo(prefix: 'INB' | 'OUT'): string {
  const now = new Date()
  const p = (n: number) => String(n).padStart(2, '0')
  const date = `${now.getFullYear()}${p(now.getMonth() + 1)}${p(now.getDate())}`
  const rand = Math.floor(Math.random() * 10000)
  return `${prefix}-${date}-${String(rand).padStart(4, '0')}`
}
