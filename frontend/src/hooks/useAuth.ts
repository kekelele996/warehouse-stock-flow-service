import { computed } from 'vue'
import { useAuthStore } from '../stores/auth'
import { ACTION_ROLES, MENU_ITEMS, type MenuItem } from '../constants/roles'

export function useAuth() {
  const store = useAuthStore()
  const menus = computed<MenuItem[]>(() => MENU_ITEMS.filter((item) => item.roles.includes(store.role as never)))
  const can = (key: keyof typeof ACTION_ROLES) => ACTION_ROLES[key].includes(store.role as never)
  const isRole = (role: string) => store.role === role
  return { store, menus, can, isRole }
}
