import request from './request'

export const getSettings = () => request.get<any, Record<string, string>>('/settings')
export const updateSettings = (data: Record<string, string>) => request.post<any, any>('/settings', data)
export const updatePassword = (old_password: string, new_password: string) =>
  request.post<any, any>('/settings/password', { old_password, new_password })

export const downloadBackup = () =>
  request.get<any, any>('/settings/backup', { responseType: 'blob' })

export const restoreBackup = (formData: FormData) =>
  request.post<any, any>('/settings/restore', formData, {
    headers: { 'Content-Type': 'multipart/form-data' },
  })

export const generate2FA = () => request.post<any, { secret: string; qr_code_url: string }>('/settings/2fa/generate')
export const enable2FA = (code: string) => request.post<any, any>('/settings/2fa/enable', { code })
export const disable2FA = () => request.post<any, any>('/settings/2fa/disable')
