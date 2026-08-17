import { ref } from 'vue'

export function usePagination(defaultPageSize = 10) {
  const page = ref(1)
  const pageSize = ref(defaultPageSize)
  const total = ref(0)
  const loading = ref(false)

  function reset() {
    page.value = 1
  }

  function setTotal(n: number) {
    total.value = n
  }

  function pageParams() {
    return { page: page.value, page_size: pageSize.value }
  }

  return { page, pageSize, total, loading, reset, setTotal, pageParams }
}
