// 轻量 toast 提示。
let seq = 0

export function toast(message: string, type: 'success' | 'error' = 'success') {
  const id = `toast-${seq++}`
  const el = document.createElement('div')
  el.id = id
  el.className = `toast toast--${type}`
  el.textContent = message
  document.body.appendChild(el)
  setTimeout(() => {
    el.remove()
  }, 2600)
}
