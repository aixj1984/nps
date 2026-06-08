<script setup lang="ts">
import { computed, ref, watch, watchEffect } from 'vue'
import SectionFields from './SectionFields.vue'
import {
  MODE_BRIDGE,
  MODE_COMMON_ONLY,
  MODE_VISITOR,
  applyWorkMode,
  cloneDocument,
  fieldDefs,
  kindLabels,
  sectionAllowsExtraFields,
  newSection,
  pruneEmptyFields,
  toWailsDocument,
  type ConfigDocument,
  type SectionKind,
} from '../configVisual'
import { ValidateConfigDocument } from '../wailsjs/go/main/NpcDeskApp.js'

const doc = defineModel<ConfigDocument>({ required: true })

defineProps<{
  path: string
}>()

const emit = defineEmits<{
  save: []
  canSaveChange: [boolean]
}>()

const workMode = ref(MODE_COMMON_ONLY)
const localErrors = ref<string[]>([])
const localWarnings = ref<string[]>([])

const sectionIndexes = computed(() =>
  doc.value.sections.map((sec, index) => ({ sec, index })).filter((row) => row.sec.kind !== 'common'),
)
const commonIndex = computed(() => doc.value.sections.findIndex((s) => s.kind === 'common'))

const modeLabel: Record<string, string> = {
  [MODE_BRIDGE]: '桥接隧道模式（TCP/UDP/HTTP/提供端 P2P）',
  [MODE_VISITOR]: '访问端模式（local_port P2P/Secret）',
  [MODE_COMMON_ONLY]: '仅连接（隧道在 Web 管理端配置）',
}

watch(
  () => doc.value.mode,
  (mode) => {
    if (mode) workMode.value = mode
  },
  { immediate: true },
)

async function validateNow() {
  const checked = normalizeValidate(await ValidateConfigDocument(toWailsDocument(doc.value)))
  localErrors.value = checked.errors ?? []
  localWarnings.value = checked.warnings ?? []
  doc.value.mode = checked.mode
}

function normalizeValidate(raw: unknown): ConfigDocument {
  const data = (raw && typeof raw === 'object' ? raw : {}) as ConfigDocument
  return {
    ...doc.value,
    mode: String((data as ConfigDocument).mode ?? doc.value.mode),
    errors: Array.isArray((data as ConfigDocument).errors) ? (data as ConfigDocument).errors : [],
    warnings: Array.isArray((data as ConfigDocument).warnings) ? (data as ConfigDocument).warnings : [],
  }
}

function setWorkMode(mode: string) {
  workMode.value = mode
  const next = applyWorkMode(doc.value, mode)
  next.mode = mode
  doc.value = next
  void validateNow()
}

function toggleSection(id: string, enabled: boolean) {
  const next = cloneDocument(doc.value)
  const sec = next.sections.find((s) => s.id === id)
  if (!sec || sec.kind === 'common') return
  sec.enabled = enabled
  doc.value = next
  void validateNow()
}

function addSection(kind: SectionKind) {
  const next = cloneDocument(doc.value)
  const count = next.sections.filter((s) => s.kind === kind).length + 1
  const sec = newSection(kind, count)
  sec.enabled = false
  next.sections.push(sec)
  if (kind === 'visitor_p2p' || kind === 'visitor_secret') {
    doc.value = applyWorkMode(next, MODE_VISITOR)
    workMode.value = MODE_VISITOR
  } else {
    doc.value = applyWorkMode(next, MODE_BRIDGE)
    workMode.value = MODE_BRIDGE
  }
  void validateNow()
}

async function triggerSave() {
  let next = applyWorkMode(doc.value, workMode.value)
  next.mode = workMode.value
  next = pruneEmptyFields(next)
  doc.value = next
  await validateNow()
  if (localErrors.value.length > 0) return
  emit('save')
}

watchEffect(() => {
  emit('canSaveChange', localErrors.value.length === 0)
})

defineExpose({ triggerSave, validateNow })
</script>

<template>
  <div class="visual-config">
    <p class="path">当前文件：<code>{{ path || '未指定' }}</code></p>
    <p class="hint">按分段管理 npc.conf；未填写的项不会写入文件；关闭分段将整段以 <code>#</code> 注释保存（不删除）</p>

    <div class="mode-bar">
      <span class="mode-label">工作模式</span>
      <button type="button" class="mode-btn" :class="{ active: workMode === MODE_COMMON_ONLY }" @click="setWorkMode(MODE_COMMON_ONLY)">
        仅连接
      </button>
      <button type="button" class="mode-btn" :class="{ active: workMode === MODE_BRIDGE }" @click="setWorkMode(MODE_BRIDGE)">
        桥接隧道
      </button>
      <button type="button" class="mode-btn" :class="{ active: workMode === MODE_VISITOR }" @click="setWorkMode(MODE_VISITOR)">
        访问端
      </button>
      <span class="mode-desc">{{ modeLabel[workMode] || doc.mode }}</span>
    </div>

    <div v-if="localErrors.length" class="msg error">
      <p v-for="(e, i) in localErrors" :key="'e' + i">{{ e }}</p>
    </div>
    <div v-if="localWarnings.length" class="msg warn">
      <p v-for="(w, i) in localWarnings" :key="'w' + i">{{ w }}</p>
    </div>

    <section v-if="commonIndex >= 0" class="card section-card">
      <div class="section-head">
        <h3>{{ kindLabels.common }}</h3>
        <span class="badge ok">始终启用</span>
      </div>
      <SectionFields
        v-model="doc.sections[commonIndex].fields"
        :defs="fieldDefs.common"
        :allow-extra="false"
      />
    </section>

    <section
      v-for="{ sec, index } in sectionIndexes"
      :key="sec.id"
      class="card section-card"
      :class="{ off: !sec.enabled }"
    >
      <div class="section-head">
        <label class="toggle">
          <input type="checkbox" :checked="sec.enabled" @change="toggleSection(sec.id, ($event.target as HTMLInputElement).checked)" />
          <strong>{{ sec.title }}</strong>
          <span class="kind">{{ kindLabels[sec.kind] || sec.kind }}</span>
        </label>
      </div>
      <SectionFields
        v-model="doc.sections[index].fields"
        :defs="fieldDefs[sec.kind] || []"
        :disabled="!sec.enabled"
        :allow-extra="sectionAllowsExtraFields(sec.kind)"
      />
    </section>

    <div class="add-row">
      <span>添加分段（默认注释，勾选开启后生效）：</span>
      <button type="button" class="bp-minimal" @click="addSection('tunnel')">+ 隧道</button>
      <button type="button" class="bp-minimal" @click="addSection('host')">+ 域名</button>
      <button type="button" class="bp-minimal" @click="addSection('health')">+ 健康检查</button>
      <button type="button" class="bp-minimal" @click="addSection('visitor_p2p')">+ P2P 访问端</button>
      <button type="button" class="bp-minimal" @click="addSection('visitor_secret')">+ Secret 访问端</button>
    </div>
  </div>
</template>

<style scoped>
.visual-config {
  margin-bottom: 16px;
}

.hint,
.path {
  margin: 0 0 8px;
  font-size: 13px;
  color: var(--muted);
}

.path {
  margin-bottom: 4px;
}

.mode-bar {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 8px;
  margin-bottom: 12px;
  padding: 10px 12px;
  background: #f5f8fa;
  border: 1px solid var(--border);
  border-radius: 3px;
}

.mode-label {
  font-size: 13px;
  color: var(--muted);
}

.mode-btn {
  border: 1px solid var(--border);
  background: #fff;
  color: var(--text);
  padding: 6px 10px;
  border-radius: 3px;
  font-size: 12px;
}

.mode-btn.active {
  border-color: var(--bp-blue);
  color: var(--bp-blue);
  background: rgba(19, 124, 189, 0.08);
}

.mode-desc {
  font-size: 12px;
  color: var(--muted);
}

.msg {
  margin-bottom: 12px;
  padding: 10px 12px;
  border-radius: 3px;
  font-size: 13px;
}

.msg p {
  margin: 0;
}

.msg.error {
  background: rgba(219, 55, 55, 0.1);
  color: var(--bp-red);
  border: 1px solid rgba(219, 55, 55, 0.2);
}

.msg.warn {
  background: rgba(217, 130, 43, 0.1);
  color: #a66321;
  border: 1px solid rgba(217, 130, 43, 0.2);
}

.section-card {
  margin-bottom: 12px;
}

.section-card.off {
  opacity: 0.82;
}

.section-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 12px;
}

h3 {
  margin: 0;
  font-size: 16px;
  color: var(--bp-blue);
}

.toggle {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
}

.kind {
  font-size: 12px;
  color: var(--muted);
}

.bp-minimal.danger {
  color: var(--bp-red);
}

.add-row {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 8px;
  font-size: 13px;
  color: var(--muted);
}

code {
  font-family: Consolas, monospace;
  color: var(--bp-blue);
}
</style>
