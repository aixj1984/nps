<script setup lang="ts">
import { ref } from 'vue'
import type { ConfigDocument } from '../configVisual'
import ConfigEditor from './ConfigEditor.vue'
import VisualConfigPanel from './VisualConfigPanel.vue'

const configContent = defineModel<string>('content', { required: true })
const configDocument = defineModel<ConfigDocument>('document', { required: true })

defineProps<{
  path: string
  saving?: boolean
}>()

const emit = defineEmits<{
  saveVisual: []
  saveRaw: []
  openFolder: []
}>()

const subTab = ref<'visual' | 'raw'>('visual')
const visualRef = ref<InstanceType<typeof VisualConfigPanel> | null>(null)
const visualCanSave = ref(true)

async function handleSave() {
  if (subTab.value === 'visual') {
    await visualRef.value?.triggerSave()
    return
  }
  emit('saveRaw')
}
</script>

<template>
  <div class="config-panel">
    <div class="config-toolbar">
      <nav class="sub-tabs">
        <button type="button" :class="{ active: subTab === 'visual' }" @click="subTab = 'visual'">分段配置</button>
        <button type="button" :class="{ active: subTab === 'raw' }" @click="subTab = 'raw'">INI 原文</button>
      </nav>
      <div class="toolbar-actions">
        <button type="button" class="bp-minimal" @click="emit('openFolder')">打开配置目录</button>
        <button
          v-if="subTab === 'visual'"
          type="button"
          class="bp-minimal"
          @click="visualRef?.validateNow()"
        >
          校验
        </button>
        <button
          type="button"
          class="bp-primary"
          :disabled="saving || (subTab === 'visual' && !visualCanSave)"
          @click="handleSave"
        >
          {{ saving ? '保存中…' : '保存配置' }}
        </button>
      </div>
    </div>

    <div class="config-body">
      <VisualConfigPanel
        v-show="subTab === 'visual'"
        ref="visualRef"
        v-model="configDocument"
        :path="path"
        @save="emit('saveVisual')"
        @can-save-change="visualCanSave = $event"
      />

      <ConfigEditor
        v-show="subTab === 'raw'"
        v-model="configContent"
        :path="path"
      />
    </div>
  </div>
</template>

<style scoped>
.config-panel {
  display: flex;
  flex-direction: column;
  flex: 1;
  min-height: 0;
}

.config-toolbar {
  flex-shrink: 0;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  padding: 10px 20px;
  background: var(--bg);
  border-bottom: 1px solid var(--border);
}

.config-body {
  flex: 1;
  min-height: 0;
  overflow-y: auto;
  padding: 16px 20px 0;
}

.sub-tabs {
  display: flex;
  gap: 4px;
  flex-shrink: 0;
}

.sub-tabs button {
  background: #fff;
  border: 1px solid var(--border);
  color: var(--muted);
  padding: 8px 14px;
  border-radius: 3px;
}

.sub-tabs button.active {
  color: var(--bp-blue);
  border-color: var(--bp-blue);
}

.toolbar-actions {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-shrink: 0;
}
</style>
