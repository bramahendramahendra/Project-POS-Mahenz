export interface PortalApp {
  id: number
  code: string
  name: string
  description: string
  path: string
  icon: string
}

export interface LaunchResponse {
  redirect_url: string
  expires_in_seconds: number
}
