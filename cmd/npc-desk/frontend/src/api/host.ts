import type { HostStatus } from '../types'
import { normalizeHostStatus } from '../hostStatus'

/** Poll the embedded HTTP API (same data as /api/host, reliable JSON). */
export async function fetchHostStatus(port: string): Promise<HostStatus> {
  const base = `http://127.0.0.1:${port}`
  const res = await fetch(`${base}/api/host`, { cache: 'no-store' })
  if (!res.ok) {
    throw new Error(`host api ${res.status}`)
  }
  const json = await res.json()
  return normalizeHostStatus(json)
}
