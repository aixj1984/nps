import { configdoc } from './wailsjs/go/models'

export type SectionKind =
  | 'common'
  | 'tunnel'
  | 'host'
  | 'health'
  | 'visitor_p2p'
  | 'visitor_secret'

export interface ConfigSection {
  id: string
  title: string
  name: string
  enabled: boolean
  kind: SectionKind | string
  fields: Record<string, string>
}

export interface ConfigDocument {
  sections: ConfigSection[]
  mode: string
  warnings?: string[]
  errors?: string[]
}

export interface FieldDef {
  key: string
  label: string
  type?: 'text' | 'number' | 'bool' | 'select'
  options?: { value: string; label: string }[]
  placeholder?: string
}

export const MODE_BRIDGE = 'bridge'
export const MODE_VISITOR = 'visitor'
export const MODE_COMMON_ONLY = 'common_only'

export const kindLabels: Record<string, string> = {
  common: '全局 [common]',
  tunnel: '隧道',
  host: '域名映射',
  health: '健康检查',
  visitor_p2p: 'P2P 访问端',
  visitor_secret: 'Secret 访问端',
}

const connTypeOptions = [
  { value: 'tcp', label: 'tcp' },
  { value: 'tls', label: 'tls' },
  { value: 'kcp', label: 'kcp' },
  { value: 'ws', label: 'ws' },
  { value: 'wss', label: 'wss' },
  { value: 'quic', label: 'quic' },
]

const tunnelModeOptions = [
  { value: 'tcp', label: 'tcp' },
  { value: 'udp', label: 'udp' },
  { value: 'httpProxy', label: 'httpProxy' },
  { value: 'socks5', label: 'socks5' },
  { value: 'file', label: 'file' },
  { value: 'p2p', label: 'p2p (提供端)' },
  { value: 'secret', label: 'secret (提供端)' },
]

const schemeOptions = [
  { value: 'all', label: 'all' },
  { value: 'http', label: 'http' },
  { value: 'https', label: 'https' },
]

const targetTypeOptions = [
  { value: 'tcp', label: 'tcp' },
  { value: 'udp', label: 'udp' },
  { value: 'all', label: 'all' },
]

const visitorTargetTypeOptions = [
  { value: 'tcp', label: 'tcp' },
  { value: 'udp', label: 'udp' },
]

export const fieldDefs: Record<string, FieldDef[]> = {
  common: [
    { key: 'server_addr', label: '服务器地址', placeholder: '127.0.0.1:8024' },
    { key: 'conn_type', label: '连接类型', type: 'select', options: connTypeOptions },
    { key: 'vkey', label: '验证密钥 vkey' },
    { key: 'dns_server', label: 'DNS 服务器', placeholder: '8.8.8.8' },
    { key: 'auto_reconnection', label: '自动重连', type: 'bool' },
    { key: 'local_ip', label: '本地出口 IP', placeholder: '留空则自动选择' },
    { key: 'ntp_server', label: 'NTP 服务器', placeholder: 'pool.ntp.org' },
    { key: 'ntp_interval', label: 'NTP 间隔(分钟)', type: 'number', placeholder: '5' },
    { key: 'max_conn', label: '最大连接数', type: 'number' },
    { key: 'flow_limit', label: '流量限制', type: 'number' },
    { key: 'rate_limit', label: '速率限制', type: 'number' },
    { key: 'time_limit', label: '时间限制', placeholder: '如 2026-12-31' },
    { key: 'remark', label: '客户端备注' },
    { key: 'basic_username', label: 'Basic 用户名' },
    { key: 'basic_password', label: 'Basic 密码' },
    { key: 'web_username', label: 'Web 用户名' },
    { key: 'web_password', label: 'Web 密码' },
    { key: 'crypt', label: '流量加密 crypt', type: 'bool' },
    { key: 'compress', label: '流量压缩 compress', type: 'bool' },
    { key: 'tls_enable', label: 'TLS 连接 tls_enable', type: 'bool' },
    { key: 'proxy_url', label: '代理 URL', placeholder: 'socks5://127.0.0.1:1080' },
    { key: 'pprof_addr', label: 'pprof 监听地址', placeholder: '0.0.0.0:9999' },
    { key: 'disconnect_timeout', label: '断开超时(秒)', type: 'number' },
  ],
  tunnel: [
    { key: 'mode', label: '模式 mode', type: 'select', options: tunnelModeOptions },
    { key: 'server_port', label: '服务端端口 server_port', type: 'number' },
    { key: 'server_ip', label: '绑定 IP server_ip', placeholder: '0.0.0.0' },
    { key: 'target_addr', label: '目标地址 target_addr', placeholder: '127.0.0.1:8080' },
    { key: 'target_ip', label: '目标 IP target_ip' },
    { key: 'target_port', label: '目标端口 target_port' },
    { key: 'password', label: '密钥 password' },
    { key: 'proxy_protocol', label: 'Proxy Protocol', type: 'number', placeholder: '0' },
    { key: 'socks5_proxy', label: '启用 SOCKS5', type: 'bool' },
    { key: 'http_proxy', label: '启用 HTTP 代理', type: 'bool' },
    { key: 'local_path', label: '本地路径 local_path (file)', placeholder: '/path/to/dir' },
    { key: 'strip_pre', label: '路径前缀 strip_pre (file)', placeholder: '/web/' },
    { key: 'read_only', label: '只读 read_only (file)', type: 'bool' },
    { key: 'dest_acl_mode', label: '目标 ACL 模式', type: 'number', placeholder: '0' },
    { key: 'dest_acl_rules', label: '目标 ACL 规则', placeholder: '逗号或换行分隔' },
    { key: 'multi_account', label: '多账号配置路径', placeholder: 'conf/multi_account.conf' },
  ],
  host: [
    { key: 'host', label: '域名 host', placeholder: 'app.example.com' },
    { key: 'target_addr', label: '目标地址 target_addr', placeholder: '127.0.0.1:8080' },
    { key: 'target_is_https', label: '目标为 HTTPS', type: 'bool' },
    { key: 'proxy_protocol', label: 'Proxy Protocol', type: 'number', placeholder: '0' },
    { key: 'host_change', label: 'Host 改写 host_change' },
    { key: 'scheme', label: '协议 scheme', type: 'select', options: schemeOptions },
    { key: 'location', label: '路径 location', placeholder: '/' },
    { key: 'path_rewrite', label: '路径改写 path_rewrite', placeholder: '/web1/' },
    { key: 'cert_file', label: '证书文件 cert_file', placeholder: 'conf/cert.pem' },
    { key: 'key_file', label: '私钥文件 key_file', placeholder: 'conf/key.pem' },
    { key: 'https_just_proxy', label: 'HTTPS 仅代理', type: 'bool' },
    { key: 'auto_ssl', label: '自动 SSL auto_ssl', type: 'bool' },
    { key: 'auto_https', label: '自动 HTTPS auto_https', type: 'bool' },
    { key: 'auto_cors', label: '自动 CORS auto_cors', type: 'bool' },
    { key: 'compat_mode', label: '兼容模式 compat_mode', type: 'bool' },
    { key: 'redirect_url', label: '重定向 URL redirect_url' },
    { key: 'multi_account', label: '多账号配置路径', placeholder: 'conf/multi_account.conf' },
    { key: 'header_X-Forwarded-Proto', label: '请求头 X-Forwarded-Proto', placeholder: 'https' },
    { key: 'response_Access-Control-Allow-Origin', label: '响应头 Access-Control-Allow-Origin', placeholder: '*' },
  ],
  health: [
    { key: 'health_check_type', label: '检查类型', type: 'select', options: [
      { value: 'http', label: 'http' },
      { value: 'tcp', label: 'tcp' },
    ]},
    { key: 'health_check_target', label: '检查目标', placeholder: '127.0.0.1:8080' },
    { key: 'health_http_url', label: 'HTTP 路径 health_http_url', placeholder: '/' },
    { key: 'health_check_interval', label: '检查间隔(秒)', type: 'number' },
    { key: 'health_check_timeout', label: '超时(秒)', type: 'number' },
    { key: 'health_check_max_failed', label: '最大失败次数', type: 'number' },
  ],
  visitor_p2p: [
    { key: 'local_port', label: '本地端口 local_port', type: 'number', placeholder: '2000' },
    { key: 'local_ip', label: '本地监听 IP local_ip' },
    { key: 'local_type', label: '本地类型 local_type' },
    { key: 'password', label: '密钥 password' },
    { key: 'target_addr', label: '目标地址 target_addr' },
    { key: 'target_type', label: '目标类型 target_type', type: 'select', options: targetTypeOptions },
    { key: 'fallback_secret', label: 'P2P 失败回落 Secret', type: 'bool' },
    { key: 'local_proxy', label: '本地代理 local_proxy', type: 'bool' },
  ],
  visitor_secret: [
    { key: 'local_port', label: '本地端口 local_port', type: 'number', placeholder: '2000' },
    { key: 'local_ip', label: '本地监听 IP local_ip' },
    { key: 'local_type', label: '本地类型 local_type' },
    { key: 'password', label: '密钥 password' },
    { key: 'target_addr', label: '目标地址 target_addr' },
    { key: 'target_type', label: '目标类型 target_type', type: 'select', options: visitorTargetTypeOptions },
    { key: 'local_proxy', label: '本地代理 local_proxy', type: 'bool' },
  ],
}

/** Section kinds that support arbitrary header_/response_ keys from the config file. */
export function sectionAllowsExtraFields(kind: string): boolean {
  return kind === 'host' || kind === 'tunnel'
}

export function normalizeDocument(raw: unknown): ConfigDocument {
  const data = (raw && typeof raw === 'object' ? raw : {}) as Record<string, unknown>
  const sectionsRaw = Array.isArray(data.sections ?? data.Sections) ? (data.sections ?? data.Sections) : []
  return {
    mode: String(data.mode ?? data.Mode ?? MODE_COMMON_ONLY),
    warnings: Array.isArray(data.warnings ?? data.Warnings) ? (data.warnings ?? data.Warnings).map(String) : [],
    errors: Array.isArray(data.errors ?? data.Errors) ? (data.errors ?? data.Errors).map(String) : [],
    sections: sectionsRaw.map((item, index) => {
      const row = (item && typeof item === 'object' ? item : {}) as Record<string, unknown>
      const fieldsRaw = (row.fields ?? row.Fields ?? {}) as Record<string, unknown>
      const fields: Record<string, string> = {}
      for (const [k, v] of Object.entries(fieldsRaw)) {
        fields[k] = String(v ?? '')
      }
      return {
        id: String(row.id ?? row.ID ?? `sec-${index}`),
        title: String(row.title ?? row.Title ?? ''),
        name: String(row.name ?? row.Name ?? ''),
        enabled: row.enabled === true || row.Enabled === true,
        kind: String(row.kind ?? row.Kind ?? 'tunnel'),
        fields,
      }
    }),
  }
}

export function toWailsDocument(doc: ConfigDocument): configdoc.Document {
  return configdoc.Document.createFrom({
    sections: doc.sections,
    mode: doc.mode,
    warnings: doc.warnings,
    errors: doc.errors,
  })
}

export function pruneEmptyFields(doc: ConfigDocument): ConfigDocument {
  const next = cloneDocument(doc)
  for (const sec of next.sections) {
    for (const [key, val] of Object.entries(sec.fields)) {
      if (String(val ?? '').trim() === '') {
        delete sec.fields[key]
      }
    }
  }
  return next
}

export function cloneDocument(doc: ConfigDocument): ConfigDocument {
  return {
    mode: doc.mode,
    warnings: doc.warnings ? [...doc.warnings] : [],
    errors: doc.errors ? [...doc.errors] : [],
    sections: doc.sections.map((s) => ({
      ...s,
      fields: { ...s.fields },
    })),
  }
}

export function applyWorkMode(doc: ConfigDocument, mode: string): ConfigDocument {
  const next = cloneDocument(doc)
  if (mode === MODE_VISITOR) {
    for (const sec of next.sections) {
      if (sec.kind === 'tunnel' || sec.kind === 'host' || sec.kind === 'health') {
        sec.enabled = false
      }
    }
  } else {
    for (const sec of next.sections) {
      if (sec.kind === 'visitor_p2p' || sec.kind === 'visitor_secret') {
        sec.enabled = false
      }
    }
  }
  return next
}

export function newSection(kind: SectionKind, index: number): ConfigSection {
  const nameByKind: Record<string, string> = {
    tunnel: `tcp${index}`,
    host: `web${index}`,
    health: `health${index}`,
    visitor_p2p: `p2p_${index}`,
    visitor_secret: `secret_${index}`,
  }
  const name = nameByKind[kind] ?? `sect${index}`
  const defaults: Record<string, Record<string, string>> = {
    tunnel: { mode: 'tcp', target_addr: '127.0.0.1:8080' },
    host: { scheme: 'all', location: '/' },
    health: { health_check_type: 'http', health_check_interval: '1' },
    visitor_p2p: { target_type: 'tcp' },
    visitor_secret: { target_type: 'tcp' },
  }
  return {
    id: `${name}-${Date.now()}`,
    title: `[${name}]`,
    name,
    enabled: false,
    kind,
    fields: { ...(defaults[kind] ?? {}) },
  }
}
