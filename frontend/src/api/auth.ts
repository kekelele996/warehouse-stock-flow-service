import { get, post } from '../utils/request'

export interface UserView {
  id: number
  username: string
  name: string
  role: string
  role_text: string
  owner_id: number | null
  status: string
}

export interface LoginResult {
  token: string
  user: UserView
}

export function apiLogin(username: string, password: string) {
  return post<LoginResult>('/auth/login', { username, password })
}

export function apiRegister(data: {
  username: string
  password: string
  name: string
  owner_name: string
  contact_name: string
  phone: string
}) {
  return post<UserView>('/auth/register', data)
}

export function apiMe() {
  return get<UserView>('/auth/me')
}
