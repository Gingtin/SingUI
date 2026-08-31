import request from './request'

export interface RoutingRule {
  id?: number
  tag: string
  type: string
  outbound: string
  domain?: string
  ip?: string
  protocol?: string
  port?: string
  network?: string
  rule_set?: string
  enable: boolean
  order: number
  remark?: string
}

export interface DNSSettings {
  id?: number
  local_dns: string
  remote_dns: string
  china_dns: string
  enable_fakeip: boolean
  fakeip_inet4?: string
  fakeip_inet6?: string
  strategy: string
  custom_rules?: string
}

export const getRoutingRules = () => request.get<any, RoutingRule[]>('/routing/rules')
export const createRoutingRule = (data: Partial<RoutingRule>) => request.post<any, RoutingRule>('/routing/rules', data)
export const updateRoutingRule = (id: number, data: Partial<RoutingRule>) => request.put<any, RoutingRule>(`/routing/rules/${id}`, data)
export const deleteRoutingRule = (id: number) => request.delete<any, any>(`/routing/rules/${id}`)

export const getDNSSettings = () => request.get<any, DNSSettings>('/routing/dns')
export const updateDNSSettings = (data: Partial<DNSSettings>) => request.post<any, any>('/routing/dns', data)
