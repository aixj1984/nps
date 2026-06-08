<script setup lang="ts">
import type { CaptureSource } from '../../types'

defineProps<{
  sources: CaptureSource[]
  loading: boolean
  selectedId: string | undefined
}>()

const emit = defineEmits<{
  loadScreen: []
  loadApp: []
  pick: [CaptureSource]
  deny: []
}>()
</script>

<template>
  <div class="select-step">
    <p class="hint">选择要共享到查看端的内容（与 Deskreen 一致）。</p>

    <div class="type-btns">
      <button type="button" class="type-card type-card--screen" @click="emit('loadScreen')">
        <span class="icon-wrap" aria-hidden="true">
          <svg class="icon-svg" viewBox="0 0 48 40" fill="none" xmlns="http://www.w3.org/2000/svg">
            <rect x="4" y="4" width="40" height="28" rx="2" stroke="currentColor" stroke-width="2" />
            <path d="M8 12h32M8 17h24M8 22h28" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" opacity="0.55" />
            <path d="M18 32h12" stroke="currentColor" stroke-width="2" stroke-linecap="round" />
            <path d="M22 32v4h4v-4" stroke="currentColor" stroke-width="2" stroke-linejoin="round" />
          </svg>
        </span>
        <span class="label">整个屏幕</span>
      </button>
      <button type="button" class="type-card type-card--app" @click="emit('loadApp')">
        <span class="icon-wrap" aria-hidden="true">
          <svg class="icon-svg" viewBox="0 0 48 40" fill="none" xmlns="http://www.w3.org/2000/svg">
            <rect x="8" y="6" width="32" height="26" rx="2" stroke="currentColor" stroke-width="2" opacity="0.45" />
            <rect x="4" y="10" width="32" height="24" rx="2" fill="currentColor" fill-opacity="0.12" stroke="currentColor" stroke-width="2" />
            <circle cx="12" cy="16" r="2" fill="currentColor" />
            <circle cx="18" cy="16" r="2" fill="currentColor" />
            <circle cx="24" cy="16" r="2" fill="currentColor" />
            <path d="M10 24h20" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" opacity="0.7" />
          </svg>
        </span>
        <span class="label">应用窗口</span>
      </button>
    </div>

    <p v-if="loading" class="loading">正在枚举显示器 / 窗口…</p>

    <div v-if="sources.length" class="source-grid">
      <button
        v-for="s in sources"
        :key="s.id"
        type="button"
        class="source-tile"
        :class="{ selected: selectedId === s.id }"
        @click="emit('pick', s)"
      >
        <img v-if="s.thumb" :src="s.thumb" alt="" />
        <div v-else class="no-thumb">无预览</div>
        <span class="name">{{ s.name }}</span>
      </button>
    </div>

    <p v-else-if="!loading" class="hint small">
      请先点击「整个屏幕」或「应用窗口」；选窗口前务必点「应用窗口」再点具体窗口。
    </p>
  </div>
</template>

<style scoped>
.select-step {
  padding: 16px 24px 24px;
}

.hint {
  color: var(--muted);
  font-size: 14px;
  margin: 0 0 12px;
}

.hint.small {
  text-align: center;
  padding: 24px 0;
}

.loading {
  text-align: center;
  color: var(--bp-blue);
  font-size: 14px;
  padding: 16px;
}

.type-btns {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 10px;
  margin-bottom: 16px;
}

.type-card {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 6px;
  padding: 10px 8px 12px;
  min-height: 0;
  border: 2px solid var(--border);
  border-radius: 3px;
  background: #f5f8fa;
  cursor: pointer;
  transition:
    border-color 0.15s,
    background 0.15s,
    box-shadow 0.15s;
}

.type-card--screen {
  --type-accent: var(--bp-blue);
}

.type-card--app {
  --type-accent: #0d8050;
}

.type-card:hover {
  border-color: var(--type-accent);
  background: #eef5fa;
  box-shadow: 0 1px 4px rgba(16, 22, 26, 0.08);
}

.type-card--app:hover {
  background: #edf7f1;
}

.type-card:hover .icon-wrap {
  background: rgba(16, 107, 163, 0.14);
}

.type-card--app:hover .icon-wrap {
  background: rgba(13, 128, 80, 0.14);
}

.icon-wrap {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 52px;
  height: 44px;
  border-radius: 4px;
  background: rgba(16, 107, 163, 0.1);
  color: var(--type-accent);
  transition: background 0.15s;
}

.type-card--app .icon-wrap {
  background: rgba(13, 128, 80, 0.1);
}

.icon-svg {
  width: 40px;
  height: 34px;
  display: block;
}

.label {
  font-weight: 600;
  font-size: 13px;
  color: #182026;
  line-height: 1.2;
}

.source-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(150px, 1fr));
  gap: 10px;
  max-height: 280px;
  overflow-y: auto;
}

.source-tile {
  border: 2px solid var(--border);
  border-radius: 3px;
  padding: 8px;
  cursor: pointer;
  text-align: center;
  background: #fff;
}

.source-tile.selected {
  border-color: var(--bp-blue);
  box-shadow: 0 0 0 2px rgba(16, 107, 163, 0.25);
}

.source-tile img {
  width: 100%;
  height: 84px;
  object-fit: contain;
  background: #1c2127;
  border-radius: 2px;
  pointer-events: none;
}

.no-thumb {
  height: 84px;
  display: grid;
  place-items: center;
  background: #1c2127;
  color: #8a9ba8;
  font-size: 12px;
}

.name {
  display: block;
  font-size: 11px;
  margin-top: 6px;
  color: var(--muted);
  word-break: break-all;
}
</style>
