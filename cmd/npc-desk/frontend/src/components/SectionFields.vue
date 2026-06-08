<script setup lang="ts">
import { computed, ref } from 'vue'
import type { FieldDef } from '../configVisual'

const fields = defineModel<Record<string, string>>({ required: true })

const props = defineProps<{
  defs: FieldDef[]
  disabled?: boolean
  allowExtra?: boolean
}>()

const newHeaderSuffix = ref('')
const newResponseSuffix = ref('')

const defKeySet = computed(() => new Set(props.defs.map((d) => d.key)))

const extraKeys = computed(() =>
  Object.keys(fields.value)
    .filter((k) => !defKeySet.value.has(k))
    .sort(),
)

function boolValue(key: string): boolean {
  const v = fields.value[key]
  return v === 'true' || v === '1'
}

function setBool(key: string, checked: boolean) {
  if (checked) {
    fields.value[key] = 'true'
  } else {
    delete fields.value[key]
  }
}

function extraLabel(key: string): string {
  if (key.startsWith('header_')) {
    return `请求头 ${key.slice('header_'.length)}`
  }
  if (key.startsWith('response_')) {
    return `响应头 ${key.slice('response_'.length)}`
  }
  return key
}

function addHeaderField() {
  const suffix = newHeaderSuffix.value.trim()
  if (!suffix) return
  const key = suffix.startsWith('header_') ? suffix : `header_${suffix}`
  if (!fields.value[key]) {
    fields.value[key] = ''
  }
  newHeaderSuffix.value = ''
}

function addResponseField() {
  const suffix = newResponseSuffix.value.trim()
  if (!suffix) return
  const key = suffix.startsWith('response_') ? suffix : `response_${suffix}`
  if (!fields.value[key]) {
    fields.value[key] = ''
  }
  newResponseSuffix.value = ''
}

function removeExtra(key: string) {
  delete fields.value[key]
}
</script>

<template>
  <div class="section-fields" :class="{ disabled }">
    <div class="field-grid">
      <label v-for="def in defs" :key="def.key" class="field">
        <span class="label">{{ def.label }}</span>
        <select
          v-if="def.type === 'select'"
          v-model="fields[def.key]"
          :disabled="disabled"
        >
          <option value=""></option>
          <option v-for="opt in def.options" :key="opt.value" :value="opt.value">{{ opt.label }}</option>
        </select>
        <input
          v-else-if="def.type === 'bool'"
          type="checkbox"
          :checked="boolValue(def.key)"
          :disabled="disabled"
          @change="setBool(def.key, ($event.target as HTMLInputElement).checked)"
        />
        <input
          v-else
          v-model="fields[def.key]"
          :type="def.type === 'number' ? 'number' : 'text'"
          :placeholder="def.placeholder"
          :disabled="disabled"
        />
      </label>
    </div>

    <template v-if="extraKeys.length">
      <h4 class="extra-title">其他参数</h4>
      <div class="field-grid">
        <label v-for="key in extraKeys" :key="key" class="field extra-field">
          <span class="label-row">
            <span class="label">{{ extraLabel(key) }}</span>
            <button
              v-if="!disabled"
              type="button"
              class="remove-btn"
               @click="removeExtra(key)"
            >
              移除
            </button>
          </span>
          <input v-model="fields[key]" type="text" :disabled="disabled" />
        </label>
      </div>
    </template>

    <div v-if="allowExtra && !disabled" class="add-extra">
      <div class="add-row">
        <input v-model="newHeaderSuffix" type="text" placeholder="Header 名，如 X-Real-IP" />
        <button type="button" class="bp-minimal" @click="addHeaderField">+ 请求头</button>
      </div>
      <div class="add-row">
        <input v-model="newResponseSuffix" type="text" placeholder="Response 名，如 X-Custom" />
        <button type="button" class="bp-minimal" @click="addResponseField">+ 响应头</button>
      </div>
    </div>
  </div>
</template>

<style scoped>
.section-fields.disabled {
  opacity: 0.55;
}

.field-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 10px 14px;
}

.field {
  display: flex;
  flex-direction: column;
  gap: 4px;
  font-size: 13px;
}

.label {
  color: var(--muted);
}

.label-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
}

.remove-btn {
  border: none;
  background: transparent;
  color: var(--bp-red);
  font-size: 12px;
  padding: 0;
  font-weight: 500;
}

.extra-title {
  margin: 14px 0 8px;
  font-size: 13px;
  color: var(--bp-blue);
  font-weight: 600;
}

.add-extra {
  margin-top: 12px;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.add-row {
  display: flex;
  gap: 8px;
  align-items: center;
}

.add-row input {
  flex: 1;
}

input,
select {
  border: 1px solid var(--border);
  border-radius: 3px;
  padding: 7px 8px;
  font-size: 13px;
}

input[type='checkbox'] {
  width: auto;
  align-self: flex-start;
}
</style>
