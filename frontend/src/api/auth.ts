import request from './request'

export const login = (data: { username: string; password: string; two_fa_code?: string }) =>
  request.post<any, { token: string; user: any }>('/auth/login', data)

export const logout = () => request.post('/auth/logout')

export const getUserInfo = () => request.get<any, { id: number; username: string; role: string; two_fa_enabled: boolean }>('/auth/info')
