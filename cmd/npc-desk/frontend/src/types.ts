export interface ServiceItem {
  name: string
  kind: 'tunnel' | 'host' | 'local' | 'health' | string
  mode: string
  port: string
  target: string
  active: boolean
}

export interface NpcStatus {
  running: boolean
  connected: boolean
  reconnecting: boolean
  hasFailed: boolean
  lastError: string
  configPath: string
  serverAddr: string
  connType: string
  vkey: string
  autoReconnection: boolean
  version: string
  activeForwardCount: number
  totalForwardCount: number
  recentlyForwarded: boolean
  services: ServiceItem[]
}

export interface CommonSettings {
  server_addr: string
  conn_type: string
  vkey: string
  dns_server: string
  auto_reconnection: boolean
  local_ip: string
  ntp_server: string
  ntp_interval: number
  max_conn: number
  flow_limit: number
  rate_limit: number
  crypt: boolean
  compress: boolean
  tls_enable: boolean
  proxy_url: string
  disconnect_timeout: number
  web_username: string
  web_password: string
  basic_username: string
  basic_password: string
}
