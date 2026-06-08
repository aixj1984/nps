<script setup lang="ts">
const emit = defineEmits<{ disconnect: [] }>()
</script>

<template>
  <div class="sharing">
    <div class="pulse-icon">📡</div>
    <h2>正在共享屏幕</h2>
    <p>查看端浏览器中应显示实时桌面画面（WebRTC）。</p>
    <p class="hint">断开连接后可从步骤 1 连接新设备。</p>
    <button class="bp-danger" @click="emit('disconnect')">断开查看端</button>
  </div>
</template>

<style scoped>
.sharing {
  padding: 40px 20px;
  text-align: center;
}

.pulse-icon {
  font-size: 3rem;
  margin-bottom: 16px;
  animation: pulse 2s ease-in-out infinite;
}

@keyframes pulse {
  50% {
    opacity: 0.5;
    transform: scale(0.95);
  }
}

.hint {
  color: var(--muted);
  font-size: 14px;
  margin-bottom: 24px;
}
</style>
