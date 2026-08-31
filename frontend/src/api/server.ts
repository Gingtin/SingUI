import request from './request'

export interface SystemStatus {
  cpu_percent: number
  mem_total: number
  mem_used: number
  mem_percent: number
  disk_total: number
  disk_used: number
  disk_percent: number
  net_upload_rate: number
  net_download_rate: number
  net_total_sent: number
  net_total_recv: number
  uptime: number
  os: string
  platform: string
  public_ip?: string
  kernel?: string
  arch?: string
  load_avg?: number[]
}

export interface CoreStatus {
  is_running: boolean
  pid: number
  start_time: string
  last_error: string
  version: string
  memory_used?: number
}

export const getSystemStatus = () => request.get<any, SystemStatus>('/server/status')
export const getCoreStatus = () => request.get<any, CoreStatus>('/server/core-status')
export const startCore = () => request.post<any, any>('/server/core/start')
export const stopCore = () => request.post<any, any>('/server/core/stop')
export const restartCore = () => request.post<any, any>('/server/core/restart')
export const getLogs = () => request.get<any, string[]>('/server/logs')
export const getActiveConnections = () => request.get<any, any>('/server/connections')
export const killConnection = (id: string) => request.delete<any, any>(`/server/connections/${id}`)
export const killAllConnections = () => request.delete<any, any>('/server/connections')
export const checkConfig = (config: string) => request.post<any, any>('/server/config/check', { config })
export const saveConfig = (config: string) => request.post<any, any>('/server/config', { config })
