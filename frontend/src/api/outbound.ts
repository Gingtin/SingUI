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
  username?: string
  password?: string
  uuid?: string
  flow?: string
  security?: string
  sni?: string
  fingerprint?: string
  alterId?: number
  security_method?: string
  obfs_type?: string
  obfs_password?: string
  up_mbps?: number
  down_mbps?: number
  method?: string
  private_key?: string
  peer_public_key?: string
  pre_shared_key?: string
  reserved?: string
  mtu?: number
  local_address_ipv4?: string
  local_address_ipv6?: string
  detour?: string
  outbounds?: string[]
  url?: string
  interval?: string
}

export const getOutbounds = () => request.get<any, Outbound[]>('/outbounds')
export const createOutbound = (data: Partial<Outbound>) => request.post<any, Outbound>('/outbounds', data)
export const updateOutbound = (id: number, data: Partial<Outbound>) => request.put<any, Outbound>(`/outbounds/${id}`, data)
export const deleteOutbound = (id: number) => request.delete<any, any>(`/outbounds/${id}`)
export const pingOutbound = (id: number) => request.post<any, { latency: number }>(`/outbounds/${id}/ping`)
export const getRawConfig = () => request.get<any, { config: string }>('/server/config')
