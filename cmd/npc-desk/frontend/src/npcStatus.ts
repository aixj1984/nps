import type { NpcStatus } from './types'

function asRecord(raw: unknown): Record<string, unknown> {
  return raw && typeof raw === 'object' ? (raw as Record<string, unknown>) : {}
}

function toCount(raw: unknown): number {
  if (typeof raw === 'number' && !Number.isNaN(raw)) {
    return raw
  }
  if (typeof raw === 'string' && raw.trim() !== '') {
    const n = Number(raw)
    return Number.isNaN(n) ? 0 : n
  }
  return 0
}

function toBool(raw: unknown): boolean {
  return raw === true || raw === 'true' || raw === 1 || raw === '1'
}

function pickCount(data: Record<string, unknown>, ...keys: string[]): number {
  for (const key of keys) {
    if (key in data) {
      return toCount(data[key])
    }
  }
  return 0
}

function pickBool(data: Record<string, unknown>, ...keys: string[]): boolean {
  for (const key of keys) {
    if (key in data) {
      return toBool(data[key])
    }
  }
  return false
}

export function normalizeNpcStatus(raw: unknown): NpcStatus {
  const data = asRecord(raw)
  const servicesRaw = Array.isArray(data.services ?? data.Services) ? (data.services ?? data.Services) : []
  return {
    running: toBool(data.running ?? data.Running),
    connected: toBool(data.connected ?? data.Connected),
    reconnecting: toBool(data.reconnecting ?? data.Reconnecting),
    hasFailed: toBool(data.hasFailed ?? data.HasFailed),
    lastError: String(data.lastError ?? data.LastError ?? ''),
    configPath: String(data.configPath ?? data.ConfigPath ?? ''),
    serverAddr: String(data.serverAddr ?? data.ServerAddr ?? ''),
    connType: String(data.connType ?? data.ConnType ?? ''),
    vkey: String(data.vkey ?? data.VKey ?? ''),
    autoReconnection: toBool(data.autoReconnection ?? data.AutoReconnection),
    version: String(data.version ?? data.Version ?? ''),
    activeForwardCount: pickCount(data, 'activeForwardCount', 'ActiveForwardCount'),
    totalForwardCount: pickCount(data, 'totalForwardCount', 'TotalForwardCount'),
    recentlyForwarded: pickBool(data, 'recentlyForwarded', 'RecentlyForwarded'),
    services: servicesRaw.map((item) => {
      const row = asRecord(item)
      return {
        name: String(row.name ?? row.Name ?? ''),
        kind: String(row.kind ?? row.Kind ?? ''),
        mode: String(row.mode ?? row.Mode ?? ''),
        port: String(row.port ?? row.Port ?? ''),
        target: String(row.target ?? row.Target ?? ''),
        active: toBool(row.active ?? row.Active),
      }
    }),
  }
}
