<script setup lang="ts">
const props = defineProps<{
  phase: string
}>()

const steps = [
  { id: 'connect', num: 1, title: '连接', desc: '扫描二维码或打开链接' },
  { id: 'select', num: 2, title: '选择', desc: '屏幕或应用窗口' },
  { id: 'confirm', num: 3, title: '确认', desc: '开始共享' },
]

function state(id: string) {
  const order = ['connect', 'select', 'confirm', 'sharing']
  const cur = order.indexOf(props.phase)
  const idx = order.indexOf(id)
  if (props.phase === 'sharing' && id === 'confirm') return 'done'
  if (idx < cur) return 'done'
  if (idx === cur) return 'active'
  return ''
}
</script>

<template>
  <nav class="stepper">
    <div v-for="s in steps" :key="s.id" class="step" :class="state(s.id)">
      <div class="circle">{{ s.num }}</div>
      <div>
        <div class="title">{{ s.title }}</div>
        <div class="desc">{{ s.desc }}</div>
      </div>
    </div>
    <div v-if="phase === 'sharing'" class="step active">
      <div class="circle">●</div>
      <div>
        <div class="title">共享中</div>
        <div class="desc">正在向查看端投屏</div>
      </div>
    </div>
  </nav>
</template>

<style scoped>
.stepper {
  display: flex;
  gap: 8px;
  padding: 16px 20px;
  background: #fff;
  border-bottom: 1px solid var(--border);
}

.step {
  flex: 1;
  display: flex;
  gap: 10px;
  padding: 10px 12px;
  border-radius: 3px;
  opacity: 0.55;
}

.step.active {
  opacity: 1;
  background: rgba(16, 107, 163, 0.08);
  box-shadow: inset 0 0 0 1px rgba(16, 107, 163, 0.25);
}

.step.done {
  opacity: 1;
}

.step.done .circle {
  background: var(--bp-green);
}

.circle {
  width: 28px;
  height: 28px;
  border-radius: 50%;
  background: var(--muted);
  color: #fff;
  display: grid;
  place-items: center;
  font-size: 13px;
  font-weight: 700;
  flex-shrink: 0;
}

.step.active .circle {
  background: var(--bp-blue);
}

.title {
  font-weight: 600;
  font-size: 14px;
}

.desc {
  font-size: 12px;
  color: var(--muted);
}
</style>
