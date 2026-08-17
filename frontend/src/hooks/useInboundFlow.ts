import { computed } from 'vue'
import type { InboundOrderView } from '../api/inbound'

export interface FlowStep {
  key: string
  label: string
}

export const INBOUND_STEPS: FlowStep[] = [
  { key: 'Pending', label: '待收货' },
  { key: 'Received', label: '已收货' },
  { key: 'QCInProgress', label: '质检' },
  { key: 'Shelved', label: '上架' },
  { key: 'Completed', label: '完成' },
]

export const OUTBOUND_STEPS: FlowStep[] = [
  { key: 'Pending', label: '待拣货' },
  { key: 'Picking', label: '拣货' },
  { key: 'Checking', label: '复核' },
  { key: 'Packing', label: '打包' },
  { key: 'Shipped', label: '发货' },
  { key: 'Completed', label: '完成' },
]

export function useInboundFlow(order: () => InboundOrderView | null) {
  const current = computed(() => {
    const o = order()
    if (!o) return 0
    const idx = INBOUND_STEPS.findIndex((s) => s.key === o.status)
    return idx < 0 ? 0 : idx
  })
  const canReceive = computed(() => order()?.status === 'Pending')
  const canQC = computed(() => order()?.status === 'Received')
  const canShelve = computed(() => order()?.status === 'QCInProgress')
  const canComplete = computed(() => order()?.status === 'Shelved')
  return { current, canReceive, canQC, canShelve, canComplete }
}

export function useOutboundFlow(order: () => { status?: string } | null) {
  const status = computed(() => order()?.status || '')
  const canPicking = computed(() => status.value === 'Pending')
  const canChecking = computed(() => status.value === 'Picking')
  const canPacking = computed(() => status.value === 'Checking')
  const canShip = computed(() => status.value === 'Packing')
  const canComplete = computed(() => status.value === 'Shipped')
  const currentStep = computed(() => {
    const idx = OUTBOUND_STEPS.findIndex((s) => s.key === status.value)
    return idx < 0 ? 0 : idx
  })
  return { canPicking, canChecking, canPacking, canShip, canComplete, currentStep }
}
