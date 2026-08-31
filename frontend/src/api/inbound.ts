import request from './request'

export interface Client {
  id?: number
  inbound_id?: number
  uuid?: string
  password?: string
  email: string
  flow?: string
  sub_token?: string
  up?: number
  down?: number
  total?: number
  expiry_time?: number
  enable?: boolean
  limit_ip?: number
  reset_day?: number
}

export interface Inbound {
  id?: number
  user_id?: number
  tag: string
  protocol: string
  port: number
  listen: string
  network: string
  security: string
  settings: string
  stream_settings: string
  sniffing: string
  enable: boolean
  remark: string
  clients?: Client[]
}

export const getInbounds = () => request.get<any, Inbound[]>('/inbounds')
export const getInbound = (id: number) => request.get<any, Inbound>(`/inbounds/${id}`)
export const createInbound = (data: Partial<Inbound>) => request.post<any, Inbound>('/inbounds', data)
export const updateInbound = (id: number, data: Partial<Inbound>) => request.put<any, Inbound>(`/inbounds/${id}`, data)
export const deleteInbound = (id: number) => request.delete<any, any>(`/inbounds/${id}`)

export const batchDeleteInbounds = (ids: number[]) => request.post<any, any>('/inbounds/batch-delete', { ids })
export const batchToggleInbounds = (ids: number[], enable: boolean) => request.post<any, any>('/inbounds/batch-toggle', { ids, enable })

export const addClient = (inboundId: number, data: Partial<Client>) => request.post<any, Client>(`/inbounds/${inboundId}/clients`, data)
export const updateClient = (inboundId: number, clientId: number, data: Partial<Client>) => request.put<any, Client>(`/inbounds/${inboundId}/clients/${clientId}`, data)
export const deleteClient = (inboundId: number, clientId: number) => request.delete<any, any>(`/inbounds/${inboundId}/clients/${clientId}`)
export const batchDeleteClients = (ids: number[]) => request.post<any, any>('/inbounds/batch-delete-clients', { ids })
export const resetClientTraffic = (inboundId: number, clientId: number) => request.post<any, any>(`/inbounds/${inboundId}/clients/${clientId}/reset`)

export const resetAllTraffic = () => request.post<any, any>('/inbounds/reset-all')
export const getRealityKeypair = () => request.get<any, { private_key: string; public_key: string; short_id: string }>('/inbounds/reality-keypair')
export const getRandomUUID = () => request.get<any, { uuid: string }>('/inbounds/random-uuid')
