<script setup lang="ts">
defineProps<{
  lanAvailable: boolean
  port: string
  bindIp: string
}>()
</script>

<template>
  <header class="host-header">
    <div class="brand">
      <div class="logo">D</div>
      <div>
        <h1>Deskreen</h1>
        <p class="tag">Community Edition · Wails</p>
      </div>
    </div>
    <div class="meta">
      <span class="badge" :class="lanAvailable ? 'ok' : 'warn'">
        {{ lanAvailable ? 'WiFi / LAN' : '无局域网' }}
      </span>
      <span class="port">{{ bindIp }}:{{ port }}</span>
    </div>
  </header>
</template>

<style scoped>
.host-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 20px;
  background: linear-gradient(180deg, #137cbd, #106ba3);
  color: #fff;
  box-shadow: 0 1px 0 rgba(16, 22, 26, 0.2);
}

.brand {
  display: flex;
  align-items: center;
  gap: 12px;
}

.logo {
  width: 40px;
  height: 40px;
  border-radius: 8px;
  background: rgba(255, 255, 255, 0.2);
  display: grid;
  place-items: center;
  font-weight: 800;
  font-size: 1.25rem;
}

h1 {
  margin: 0;
  font-size: 1.35rem;
  font-weight: 600;
}

.tag {
  margin: 2px 0 0;
  font-size: 0.75rem;
  opacity: 0.85;
}

.meta {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
  gap: 6px;
}

.port {
  font-size: 0.8rem;
  opacity: 0.9;
  font-family: Consolas, monospace;
}

.host-header .badge.ok {
  background: rgba(255, 255, 255, 0.2);
  color: #fff;
}

.host-header .badge.warn {
  background: rgba(219, 55, 55, 0.35);
  color: #fff;
}
</style>
