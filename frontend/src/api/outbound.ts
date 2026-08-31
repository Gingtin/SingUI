import request from './request'

export interface Outbound {
  id?: number
  tag: string
  type: string // direct, block, dns, warp, wireguard, socks, http, custom
  server?: string
  port?: number
  settings?: string
  enable: boolean
  remark?: string
}

export const getOutbounds = () => request.get<any, Outbound[]>('/outbounds')
export const createOutbound = (data: Partial<Outbound>) => request.post<any, Outbound>('/outbounds', data)
export const updateOutbound = (id: number, data: Partial<Outbound>) => request.put<any, Outbound>(`/outbounds/${id}`, data)
export const deleteOutbound = (id: number) => request.delete<any, any>(`/outbounds/${id}`)
export const getRawConfig = () => request.get<any, { config: string }>('/server/config')
