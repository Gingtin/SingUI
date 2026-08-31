import request from './request'

export interface VersionInfo {
  panel: {
    current_version: string
    latest_version: string
    has_update: boolean
    release_notes: string
    release_url: string
  }
  core: {
    current_version: string
    latest_version: string
    has_update: boolean
    available_versions: string[]
  }
  geo: {
    last_updated: string
    status: string
  }
}

export const checkVersions = () => request.get<any, VersionInfo>('/version/check')
export const updateCore = (version: string) =>
  request.post<any, any>('/version/update-core', { version })
export const updateGeo = () => request.post<any, any>('/version/update-geo')
