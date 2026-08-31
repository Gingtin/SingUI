import request from './request'

export const getSettings = () => request.get<any, Record<string, string>>('/settings')
export const updateSettings = (data: Record<string, string>) => request.post<any, any>('/settings', data)
export const updatePassword = (data: { old_password: string; new_password: string }) =>
  request.post<any, any>('/settings/password', data)
