<template>
  <a
    :href="item.url"
    target="_blank"
    rel="noopener noreferrer"
    class="group relative flex flex-col justify-between rounded-xl border border-gray-200 bg-white p-5 shadow-sm transition-all duration-300 hover:-translate-y-1 hover:shadow-lg dark:border-gray-700 dark:bg-gray-800"
  >
    <div class="flex items-start justify-between">
      <div class="flex items-center gap-3">
        <div
          class="flex h-10 w-10 shrink-0 items-center justify-center rounded-lg bg-indigo-50 text-lg font-bold text-indigo-600 dark:bg-indigo-900/30 dark:text-indigo-400"
        >
          <span v-if="item.icon">{{ item.icon }}</span>
          <span v-else>{{ item.title.charAt(0).toUpperCase() }}</span>
        </div>
        <div>
          <h3 class="text-base font-semibold text-gray-900 dark:text-gray-100">
            {{ item.title }}
          </h3>
          <span class="text-xs text-gray-400">{{ item.url }}</span>
        </div>
      </div>
      <span
        v-if="online !== undefined"
        class="h-2.5 w-2.5 rounded-full"
        :class="online ? 'bg-green-500' : 'bg-red-500'"
        :title="online ? '在线' : '离线'"
      />
    </div>
    <p
      v-if="item.description"
      class="mt-3 line-clamp-2 text-sm text-gray-500 dark:text-gray-400"
    >
      {{ item.description }}
    </p>
    <div class="mt-3 flex gap-2">
      <span
        v-for="tag in tags"
        :key="tag"
        class="inline-flex items-center rounded-md bg-gray-100 px-2 py-0.5 text-xs font-medium text-gray-600 dark:bg-gray-700 dark:text-gray-300"
      >
        {{ tag }}
      </span>
    </div>
  </a>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  item: { type: Object, required: true },
  online: { type: Boolean, default: undefined },
})

const tags = computed(() => {
  const t = []
  const url = props.item.url.toLowerCase()
  if (url.includes('prod') || url.includes('生产')) t.push('Prod')
  else if (url.includes('test') || url.includes('测试')) t.push('Test')
  else if (url.includes('dev') || url.includes('开发')) t.push('Dev')
  return t
})
</script>
