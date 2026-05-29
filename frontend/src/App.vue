<template>
  <div class="flex h-screen flex-col overflow-hidden">
    <!-- Header -->
    <header class="flex items-center justify-between border-b border-gray-200 bg-white px-6 py-3 dark:border-gray-700 dark:bg-gray-800">
      <div class="flex items-center gap-3">
        <div class="flex h-8 w-8 items-center justify-center rounded-lg bg-indigo-600 text-white font-bold">N</div>
        <h1 class="text-lg font-bold text-gray-900 dark:text-white">内部导航</h1>
      </div>
      <div class="flex items-center gap-3">
        <div class="relative hidden sm:block">
          <input
            ref="searchInput"
            v-model="keyword"
            placeholder="搜索系统、URL 或描述..."
            class="w-72 rounded-lg border border-gray-300 bg-gray-50 py-1.5 pl-9 pr-4 text-sm text-gray-900 outline-none focus:border-indigo-500 focus:ring-1 focus:ring-indigo-500 dark:border-gray-600 dark:bg-gray-700 dark:text-white"
            @keydown.esc="keyword = ''"
          />
          <svg class="absolute left-2.5 top-2 h-4 w-4 text-gray-400" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
          </svg>
          <span class="absolute right-2 top-1.5 rounded border border-gray-200 bg-white px-1 text-[10px] text-gray-400 dark:border-gray-600 dark:bg-gray-700">Ctrl K</span>
        </div>
        <button
          class="rounded-lg p-2 text-gray-500 hover:bg-gray-100 dark:text-gray-300 dark:hover:bg-gray-700"
          title="切换主题"
          @click="toggleDark"
        >
          <svg v-if="isDark" class="h-5 w-5" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" d="M12 3v1m0 16v1m9-9h-1M4 12H3m15.364 6.364l-.707-.707M6.343 6.343l-.707-.707m12.728 0l-.707.707M6.343 17.657l-.707.707M16 12a4 4 0 11-8 0 4 4 0 018 0z" />
          </svg>
          <svg v-else class="h-5 w-5" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" d="M20.354 15.354A9 9 0 018.646 3.646 9.003 9.003 0 0012 21a9.003 9.003 0 008.354-5.646z" />
          </svg>
        </button>
        <button
          class="rounded-lg p-2 text-gray-500 hover:bg-gray-100 dark:text-gray-300 dark:hover:bg-gray-700"
          title="管理后台"
          @click="showAdmin = true"
        >
          <svg class="h-5 w-5" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z" />
            <path stroke-linecap="round" stroke-linejoin="round" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
          </svg>
        </button>
      </div>
    </header>

    <!-- Body -->
    <div class="flex flex-1 overflow-hidden">
      <!-- Sidebar -->
      <aside class="hidden w-56 flex-col border-r border-gray-200 bg-white dark:border-gray-700 dark:bg-gray-800 md:flex">
        <div class="p-4">
          <h2 class="mb-2 text-xs font-semibold uppercase tracking-wider text-gray-500 dark:text-gray-400">分类</h2>
          <nav class="space-y-1">
            <a
              v-for="cat in navTree"
              :key="cat.category_id"
              :href="`#cat-${cat.category_id}`"
              class="block rounded-lg px-3 py-2 text-sm font-medium text-gray-700 hover:bg-gray-100 dark:text-gray-300 dark:hover:bg-gray-700"
              :class="{ 'bg-indigo-50 text-indigo-700 dark:bg-indigo-900/30 dark:text-indigo-300': activeCat === cat.category_id }"
              @click="activeCat = cat.category_id"
            >
              {{ cat.category_name }}
              <span class="ml-1 text-xs text-gray-400">{{ cat.items?.length || 0 }}</span>
            </a>
          </nav>
        </div>
      </aside>

      <!-- Main -->
      <main ref="mainRef" class="flex-1 overflow-y-auto p-6">
        <div v-if="filteredTree.length === 0" class="flex h-full items-center justify-center text-gray-400">
          暂无数据
        </div>
        <div v-for="cat in filteredTree" :key="cat.category_id" :id="`cat-${cat.category_id}`" class="mb-10">
          <h2 class="mb-4 flex items-center gap-2 text-lg font-bold text-gray-800 dark:text-gray-100">
            <span class="h-5 w-1 rounded-full bg-indigo-600" />
            {{ cat.category_name }}
            <span class="text-xs font-normal text-gray-400">({{ cat.items?.length || 0 }})</span>
          </h2>
          <div class="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
            <NavCard
              v-for="item in cat.items"
              :key="item.id"
              :item="item"
              :online="item.online"
            />
          </div>
        </div>
      </main>
    </div>
  </div>

  <AdminModal v-model="showAdmin" :categories="categories" :all-navs="allNavs" @refresh="loadData" />
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted, watch } from 'vue'
import NavCard from './components/NavCard.vue'
import AdminModal from './components/AdminModal.vue'

const navTree = ref([])
const categories = ref([])
const allNavs = ref([])
const keyword = ref('')
const isDark = ref(false)
const showAdmin = ref(false)
const activeCat = ref(null)
const mainRef = ref(null)
const searchInput = ref(null)

const filteredTree = computed(() => {
  if (!keyword.value.trim()) return navTree.value
  const k = keyword.value.toLowerCase()
  return navTree.value.map(cat => ({
    ...cat,
    items: cat.items.filter(item =>
      item.title.toLowerCase().includes(k) ||
      item.url.toLowerCase().includes(k) ||
      (item.description && item.description.toLowerCase().includes(k))
    )
  })).filter(cat => cat.items.length > 0)
})

function toggleDark() {
  isDark.value = !isDark.value
  document.documentElement.classList.toggle('dark', isDark.value)
  localStorage.setItem('theme', isDark.value ? 'dark' : 'light')
}

function initTheme() {
  const saved = localStorage.getItem('theme')
  const prefers = window.matchMedia('(prefers-color-scheme: dark)').matches
  isDark.value = saved === 'dark' || (!saved && prefers)
  document.documentElement.classList.toggle('dark', isDark.value)
}

async function loadData() {
  const res = await fetch('/api/v1/navs')
  navTree.value = await res.json()
  const cRes = await fetch('/api/v1/categories')
  categories.value = await cRes.json()
  allNavs.value = navTree.value.flatMap(c => c.items.map(i => ({ ...i, category_id: c.category_id })))
}

function onKeydown(e) {
  if (e.ctrlKey && (e.key === 'k' || e.key === 'K')) {
    e.preventDefault()
    searchInput.value?.focus()
  }
}

function onScroll() {
  if (!mainRef.value) return
  const tops = navTree.value.map(cat => {
    const el = document.getElementById(`cat-${cat.category_id}`)
    return el ? { id: cat.category_id, top: el.offsetTop - mainRef.value.scrollTop } : null
  }).filter(Boolean)
  const active = tops.filter(t => t.top <= 20).pop()
  if (active) activeCat.value = active.id
}

onMounted(() => {
  initTheme()
  loadData()
  window.addEventListener('keydown', onKeydown)
  mainRef.value?.addEventListener('scroll', onScroll)
})

onUnmounted(() => {
  window.removeEventListener('keydown', onKeydown)
  mainRef.value?.removeEventListener('scroll', onScroll)
})
</script>
