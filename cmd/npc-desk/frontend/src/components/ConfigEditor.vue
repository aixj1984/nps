<script setup lang="ts">
const model = defineModel<string>({ required: true })

defineProps<{
  path: string
}>()
</script>

<template>
  <div class="card config-editor">
    <p class="path">当前文件：<code>{{ path || '未指定' }}</code></p>
    <p class="hint">支持 npc.conf 全部段落：[common]、隧道、域名、健康检查、P2P/Secret 等</p>
    <textarea v-model="model" spellcheck="false" placeholder="[common]&#10;server_addr=..."></textarea>
  </div>
</template>

<style scoped>
.config-editor {
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

textarea {
  width: 100%;
  min-height: 320px;
  resize: vertical;
  border: 1px solid var(--border);
  border-radius: 3px;
  padding: 12px;
  font-family: Consolas, 'Courier New', monospace;
  font-size: 13px;
  line-height: 1.5;
  background: #fbfcfd;
}

code {
  font-family: Consolas, monospace;
  color: var(--bp-blue);
}
</style>
