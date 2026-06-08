import type { DeviceInfo, HostStatus, RoomState } from './types'

function asRecord(v: unknown): Record<string, unknown> {
  return v && typeof v === 'object' ? (v as Record<string, unknown>) : {}
}

function str(v: unknown, fallback = ''): string {
  return v != null ? String(v) : fallback
}

function num(v: unknown): number {
  if (typeof v === 'number' && !Number.isNaN(v)) return v
  if (typeof v === 'string') return parseInt(v, 10) || 0
  return 0
}

function normalizeDevice(raw: unknown): DeviceInfo | undefined {
  if (!raw || typeof raw !== 'object') return undefined
  const d = raw as Record<string, unknown>
  const dev: DeviceInfo = {
    ip: str(d.ip ?? d.IP),
    userAgent: str(d.userAgent ?? d.UserAgent),
    os: str(d.os ?? d.OS),
    deviceType: str(d.deviceType ?? d.DeviceType),
    browser: str(d.browser ?? d.Browser),
    screenWidth: num(d.screenWidth ?? d.ScreenWidth),
    screenHeight: num(d.screenHeight ?? d.ScreenHeight),
  }
  if (!dev.ip && !dev.browser && !dev.os) return undefined
  return dev
}

function normalizeRoom(raw: unknown): RoomState {
  const room = asRecord(raw)
  const src = asRecord(room.selectedSource ?? room.SelectedSource)
  const pending = normalizeDevice(room.pendingDevice ?? room.PendingDevice)
  return {
    roomId: str(room.roomId ?? room.RoomID),
    phase: str(room.phase ?? room.Phase, 'connect'),
    sharingType: str(room.sharingType ?? room.SharingType),
    selectedSource: {
      id: str(src.id),
      name: str(src.name),
      type: str(src.type),
    },
    pendingDevice: pending,
    viewerOnline: !!(room.viewerOnline ?? room.ViewerOnline),
    connectedIp: str(room.connectedIp ?? room.ConnectedIP),
  }
}

/** Normalize Wails GetStatus / EventsEmit payloads (camelCase or PascalCase). */
export function normalizeHostStatus(raw: unknown): HostStatus {
  const r = asRecord(raw)
  return {
    viewerUrl: str(r.viewerUrl ?? r.ViewerURL),
    lanAvailable: !!(r.lanAvailable ?? r.LanAvailable),
    viewerOnline: !!(r.viewerOnline ?? r.ViewerOnline),
    room: normalizeRoom(r.room ?? r.Room),
  }
}
