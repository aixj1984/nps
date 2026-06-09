<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref } from 'vue'
import { EventsOn } from './wailsjs/runtime/runtime'
import NpcHeader from './components/NpcHeader.vue'
import ConnectionPanel from './components/ConnectionPanel.vue'
import ServicesPanel from './components/ServicesPanel.vue'
import ConfigPanel from './components/ConfigPanel.vue'
import {
  cloneDocument,
  documentsEqual,
  normalizeDocument,
  toWailsDocument,
  type ConfigDocument,
} from './configVisual'
import LogsPanel from './components/LogsPanel.vue'
import type { NpcStatus } from './types'
import { normalizeNpcStatus } from './npcStatus'
import {
  GetStatus,
  GetConfigContent,
  GetConfigDocument,
  GetConfigPath,
  SaveConfig,
  SaveConfigDocument,
  PickConfigFile,
  StartClient,
  StopClient,
  GetLogs,
  ClearLogs,
  OpenConfigFolder,
} from './wailsjs/go/main/NpcDeskApp.js'

const status = ref<NpcStatus | null>(null)
const configContent = ref('')
const configDocument = ref<ConfigDocument>({ sections: [], mode: 'common_only' })
const configBaselineContent = ref('')
const configBaselineDocument = ref<ConfigDocument>({ sections: [], mode: 'common_only' })
const configPath = ref('')
const logs = ref('')
const error = ref('')
const busy = ref(false)
const saving = ref(false)
const activeTab = ref<'dashboard' | 'config' | 'logs'>('dashboard')

const configDirty = computed(
  () =>
    configContent.value !== configBaselineContent.value
    || !documentsEqual(configDocument.value, configBaselineDocument.value),
)

function captureConfigBaseline() {
  configBaselineContent.value = configContent.value
  configBaselineDocument.value = cloneDocument(configDocument.value)
}

function cancelConfigEdit() {
  configContent.value = configBaselineContent.value
  configDocument.value = cloneDocument(configBaselineDocument.value)
  error.value = ''
}

async function refresh() {
  try {
    status.value = normalizeNpcStatus(await GetStatus())
    error.value = ''
  } catch (e) {
    error.value = e instanceof Error ? e.message : String(e)
  }
}

async function loadConfig() {
  configContent.value = await GetConfigContent()
  configPath.value = await GetConfigPath()
  configDocument.value = normalizeDocument(await GetConfigDocument())
  captureConfigBaseline()
}

async function loadLogs() {
  const text = await GetLogs()
  if (text) {
    logs.value = text
  }
}

async function clearLogs() {
  await ClearLogs()
  logs.value = ''
}

async function saveConfigRaw() {
  saving.value = true
  try {
    await SaveConfig(configContent.value)
    await refresh()
    await loadConfig()
  } catch (e) {
    error.value = e instanceof Error ? e.message : String(e)
  } finally {
    saving.value = false
  }
}

async function saveConfigVisual() {
  saving.value = true
  try {
    await SaveConfigDocument(toWailsDocument(configDocument.value))
    await refresh()
    await loadConfig()
  } catch (e) {
    error.value = e instanceof Error ? e.message : String(e)
  } finally {
    saving.value = false
  }
}

async function startClient() {
  busy.value = true
  try {
    await SaveConfigDocument(toWailsDocument(configDocument.value))
    await StartClient()
    await refresh()
    await loadLogs()
  } catch (e) {
    error.value = e instanceof Error ? e.message : String(e)
  } finally {
    busy.value = false
  }
}

async function stopClient() {
  busy.value = true
  try {
    await StopClient()
    await refresh()
    await loadLogs()
  } catch (e) {
    error.value = e instanceof Error ? e.message : String(e)
  } finally {
    busy.value = false
  }
}

async function pickConfig() {
  try {
    const path = await PickConfigFile()
    if (path) {
      configPath.value = path
      await loadConfig()
      await refresh()
    }
  } catch (e) {
    error.value = e instanceof Error ? e.message : String(e)
  }
}

async function openConfigFolder() {
  try {
    await OpenConfigFolder()
    error.value = ''
  } catch (e) {
    error.value = e instanceof Error ? e.message : String(e)
  }
}

let unlisten: (() => void) | null = null

onMounted(async () => {
  await refresh()
  await loadConfig()
  await loadLogs()
  unlisten = EventsOn('npc:update', (data: unknown) => {
    status.value = normalizeNpcStatus(data)
  })
  setInterval(() => {
    void refresh()
    if (activeTab.value === 'logs') void loadLogs()
  }, 500)
})

onUnmounted(() => unlisten?.())
</script>

<template>
  <div class="app">
    <NpcHeader :status="status" />

    <nav class="tabs">
      <button type="button" :class="{ active: activeTab === 'dashboard' }" @click="activeTab = 'dashboard'">
        总览
      </button>
      <button type="button" :class="{ active: activeTab === 'config' }" @click="activeTab = 'config'">
        配置
      </button>
      <button type="button" :class="{ active: activeTab === 'logs' }" @click="activeTab = 'logs'; loadLogs()">
        日志
      </button>
    </nav>

    <div v-if="error" class="error-bar">{{ error }}</div>

    <main class="main" :class="{ 'main--config': activeTab === 'config' }">
      <template v-if="activeTab === 'dashboard' && status">
        <ConnectionPanel
          :status="status"
          :busy="busy"
          @start="startClient"
          @stop="stopClient"
          @pick-config="pickConfig"
          @open-folder="openConfigFolder"
        />
        <ServicesPanel
          :services="status.services"
          :connected="status.connected"
          :active-forward-count="status.activeForwardCount"
          :total-forward-count="status.totalForwardCount"
          :recently-forwarded="status.recentlyForwarded"
        />
      </template>

      <template v-else-if="activeTab === 'config'">
        <ConfigPanel
          v-model:content="configContent"
          v-model:document="configDocument"
          :path="configPath"
          :saving="saving"
          :dirty="configDirty"
          @save-visual="saveConfigVisual"
          @save-raw="saveConfigRaw"
          @cancel="cancelConfigEdit"
          @open-folder="openConfigFolder"
        />
      </template>

      <template v-else>
        <LogsPanel :logs="logs" @refresh="loadLogs" @clear="clearLogs" />
      </template>
    </main>

    <footer class="footer">
      <button type="button" class="bp-minimal" @click="refresh">刷新状态</button>
      <button type="button" class="bp-minimal" @click="pickConfig">切换配置文件</button>
    </footer>
  </div>
</template>

<style scoped>
.app {
  display: flex;
  flex-direction: column;
  height: 100vh;
}

.tabs {
  display: flex;
  gap: 4px;
  padding: 10px 20px 0;
  background: var(--bg);
}

.tabs button {
  background: transparent;
  border: none;
  color: var(--muted);
  padding: 10px 16px;
  border-bottom: 2px solid transparent;
  border-radius: 0;
}

.tabs button.active {
  color: var(--bp-blue);
  border-bottom-color: var(--bp-blue);
}

.main {
  flex: 1;
  overflow-y: auto;
  padding: 16px 20px 0;
}

.main.main--config {
  overflow: hidden;
  padding: 0;
  display: flex;
  flex-direction: column;
}

.error-bar {
  background: rgba(219, 55, 55, 0.12);
  color: var(--bp-red);
  padding: 10px 20px;
  font-size: 13px;
}

.footer {
  display: flex;
  gap: 12px;
  padding: 10px 20px;
  border-top: 1px solid var(--border);
  background: #fff;
}
</style>
