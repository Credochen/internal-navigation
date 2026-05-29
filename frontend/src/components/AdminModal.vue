<template>
  <Teleport to="body">
    <div
      v-if="modelValue"
      class="fixed inset-0 z-50 flex items-center justify-center bg-black/50 backdrop-blur-sm"
      @click.self="close"
    >
      <div
        class="mx-4 w-full max-w-3xl max-h-[85vh] overflow-y-auto rounded-2xl bg-white p-6 shadow-2xl dark:bg-gray-800"
      >
        <div class="mb-4 flex items-center justify-between">
          <h2 class="text-xl font-bold text-gray-900 dark:text-white">管理后台</h2>
          <button
            class="rounded-lg p-1 text-gray-400 hover:bg-gray-100 hover:text-gray-600 dark:hover:bg-gray-700 dark:hover:text-gray-200"
            @click="close"
          >
            <svg class="h-6 w-6" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>

        <!-- Token Authentication -->
        <div v-if="!token" class="py-10 text-center">
          <div class="mb-2 text-sm text-gray-500 dark:text-gray-400">请输入管理员 Token 以继续</div>
          <form class="mx-auto flex max-w-sm gap-2" @submit.prevent="saveToken">
            <input
              v-model="tokenInput"
              type="password"
              placeholder="X-Admin-Token"
              class="flex-1 rounded-lg border border-gray-300 px-3 py-2 text-sm dark:border-gray-600 dark:bg-gray-700 dark:text-white"
            />
            <button type="submit" class="rounded-lg bg-indigo-600 px-4 py-2 text-sm font-medium text-white hover:bg-indigo-700">进入</button>
          </form>
          <div v-if="tokenError" class="mt-2 text-sm text-red-500">Token 验证失败，请重试</div>
        </div>

        <template v-else>
          <div class="mb-4 flex items-center justify-between">
            <div class="flex gap-2">
              <button
                v-for="t in ['categories','navs']"
                :key="t"
                class="rounded-lg px-4 py-2 text-sm font-medium transition-colors"
                :class="tab===t ? 'bg-indigo-600 text-white' : 'bg-gray-100 text-gray-600 hover:bg-gray-200 dark:bg-gray-700 dark:text-gray-300'"
                @click="tab=t"
              >
                {{ t==='categories' ? '分类管理' : '导航管理' }}
              </button>
            </div>
            <button class="text-xs text-gray-400 hover:text-indigo-600 dark:hover:text-indigo-400" @click="clearToken">退出登录</button>
          </div>

        <div v-if="tab==='categories'" class="space-y-4">
          <form class="flex gap-2" @submit.prevent="saveCategory">
            <input v-model="catForm.name" placeholder="分类名称" required class="flex-1 rounded-lg border border-gray-300 px-3 py-2 text-sm dark:border-gray-600 dark:bg-gray-700 dark:text-white" />
            <input v-model.number="catForm.sort_order" type="number" placeholder="排序" class="w-20 rounded-lg border border-gray-300 px-3 py-2 text-sm dark:border-gray-600 dark:bg-gray-700 dark:text-white" />
            <button type="submit" class="rounded-lg bg-indigo-600 px-4 py-2 text-sm font-medium text-white hover:bg-indigo-700">{{ catForm.id ? '更新' : '添加' }}</button>
            <button v-if="catForm.id" type="button" class="rounded-lg bg-gray-200 px-4 py-2 text-sm text-gray-700 hover:bg-gray-300 dark:bg-gray-700 dark:text-gray-200" @click="resetCat">取消</button>
          </form>
          <table class="w-full text-left text-sm">
            <thead class="bg-gray-50 text-gray-600 dark:bg-gray-700 dark:text-gray-300">
              <tr><th class="px-3 py-2">名称</th><th class="px-3 py-2">排序</th><th class="px-3 py-2 text-right">操作</th></tr>
            </thead>
            <tbody class="divide-y divide-gray-100 dark:divide-gray-700">
              <tr v-for="c in categories" :key="c.id" class="hover:bg-gray-50 dark:hover:bg-gray-700/50">
                <td class="px-3 py-2 dark:text-gray-100">{{ c.name }}</td>
                <td class="px-3 py-2 dark:text-gray-100">{{ c.sort_order }}</td>
                <td class="px-3 py-2 text-right space-x-2">
                  <button class="text-indigo-600 hover:underline dark:text-indigo-400" @click="editCat(c)">编辑</button>
                  <button class="text-red-600 hover:underline dark:text-red-400" @click="delCat(c.id)">删除</button>
                </td>
              </tr>
            </tbody>
          </table>
        </div>

        <div v-else class="space-y-4">
          <form class="grid grid-cols-1 gap-2 sm:grid-cols-2" @submit.prevent="saveNav">
            <input v-model="navForm.title" placeholder="标题" required class="rounded-lg border border-gray-300 px-3 py-2 text-sm dark:border-gray-600 dark:bg-gray-700 dark:text-white" />
            <input v-model="navForm.url" placeholder="URL" required class="rounded-lg border border-gray-300 px-3 py-2 text-sm dark:border-gray-600 dark:bg-gray-700 dark:text-white" />
            <select v-model.number="navForm.category_id" required class="rounded-lg border border-gray-300 px-3 py-2 text-sm dark:border-gray-600 dark:bg-gray-700 dark:text-white">
              <option value="0">选择分类</option>
              <option v-for="c in categories" :key="c.id" :value="c.id">{{ c.name }}</option>
            </select>
            <select v-model="navForm.icon" class="rounded-lg border border-gray-300 px-3 py-2 text-sm dark:border-gray-600 dark:bg-gray-700 dark:text-white">
              <option value="">无图标</option>
              <option value="🖥️">🖥️ 服务器</option>
              <option value="📊">📊 监控/数据</option>
              <option value="🚀">🚀 CI/CD</option>
              <option value="🗄️">🗄️ 数据库</option>
              <option value="🔧">🔧 工具</option>
              <option value="🌐">🌐 网络/Web</option>
              <option value="📁">📁 文件/存储</option>
              <option value="🔒">🔒 安全</option>
              <option value="📡">📡 消息队列</option>
              <option value="⚙️">⚙️ 配置/设置</option>
              <option value="📝">📝 日志</option>
              <option value="🔍">🔍 搜索</option>
              <option value="☁️">☁️ 云/容器</option>
              <option value="📈">📈 统计</option>
              <option value="🔥">🔥 告警</option>
              <option value="🐳">🐳 Docker</option>
              <option value="☸️">☸️ Kubernetes</option>
              <option value="🐧">🐧 Linux</option>
              <option value="🐹">🐹 Go</option>
              <option value="🐍">🐍 Python</option>
              <option value="☕">☕ Java</option>
              <option value="📦">📦 制品库</option>
              <option value="🧪">🧪 测试</option>
              <option value="📋">📋 文档</option>
              <option value="💬">💬 通讯</option>
            </select>
            <input v-model="navForm.description" placeholder="描述（可选）" class="rounded-lg border border-gray-300 px-3 py-2 text-sm dark:border-gray-600 dark:bg-gray-700 dark:text-white" />
            <input v-model.number="navForm.sort_order" type="number" placeholder="排序" class="rounded-lg border border-gray-300 px-3 py-2 text-sm dark:border-gray-600 dark:bg-gray-700 dark:text-white" />
            <div class="flex gap-2 sm:col-span-2">
              <button type="submit" class="rounded-lg bg-indigo-600 px-4 py-2 text-sm font-medium text-white hover:bg-indigo-700">{{ navForm.id ? '更新' : '添加' }}</button>
              <button v-if="navForm.id" type="button" class="rounded-lg bg-gray-200 px-4 py-2 text-sm text-gray-700 hover:bg-gray-300 dark:bg-gray-700 dark:text-gray-200" @click="resetNav">取消</button>
            </div>
          </form>
          <table class="w-full text-left text-sm">
            <thead class="bg-gray-50 text-gray-600 dark:bg-gray-700 dark:text-gray-300">
              <tr><th class="px-3 py-2">标题</th><th class="px-3 py-2">URL</th><th class="px-3 py-2">分类</th><th class="px-3 py-2 text-right">操作</th></tr>
            </thead>
            <tbody class="divide-y divide-gray-100 dark:divide-gray-700">
              <tr v-for="n in allNavs" :key="n.id" class="hover:bg-gray-50 dark:hover:bg-gray-700/50">
                <td class="px-3 py-2 dark:text-gray-100">{{ n.title }}</td>
                <td class="px-3 py-2 dark:text-gray-100 truncate max-w-xs">{{ n.url }}</td>
                <td class="px-3 py-2 dark:text-gray-100">{{ catName(n.category_id) }}</td>
                <td class="px-3 py-2 text-right space-x-2">
                  <button class="text-indigo-600 hover:underline dark:text-indigo-400" @click="editNav(n)">编辑</button>
                  <button class="text-red-600 hover:underline dark:text-red-400" @click="delNav(n.id)">删除</button>
                </td>
              </tr>
            </tbody>
          </table>
        </div>

        <div class="mt-6 flex gap-2 border-t border-gray-100 pt-4 dark:border-gray-700">
          <button class="rounded-lg bg-gray-100 px-3 py-2 text-sm text-gray-700 hover:bg-gray-200 dark:bg-gray-700 dark:text-gray-200" @click="doExport">导出 JSON</button>
          <label class="cursor-pointer rounded-lg bg-gray-100 px-3 py-2 text-sm text-gray-700 hover:bg-gray-200 dark:bg-gray-700 dark:text-gray-200">
            导入 JSON
            <input type="file" accept="application/json" class="hidden" @change="doImport" />
          </label>
        </div>
        </template>
      </div>
    </div>
  </Teleport>
</template>

<script setup>
import { ref, watch } from 'vue'

const props = defineProps({
  modelValue: Boolean,
  categories: { type: Array, default: () => [] },
  allNavs: { type: Array, default: () => [] },
})
const emit = defineEmits(['update:modelValue', 'refresh'])

const tab = ref('categories')
const token = ref(localStorage.getItem('adminToken') || '')
const tokenInput = ref('')
const tokenError = ref(false)

const api = (path, opts) => {
  const headers = { 'X-Admin-Token': token.value, ...(opts?.headers || {}) }
  return fetch('/api/v1/admin' + path, { ...opts, headers })
}

const catForm = ref({ id: 0, name: '', sort_order: 0 })
const navForm = ref({ id: 0, title: '', url: '', category_id: 0, icon: '', description: '', sort_order: 0 })

function close() {
  emit('update:modelValue', false)
}

async function saveToken() {
  tokenError.value = false
  const test = await fetch('/api/v1/admin/export', { headers: { 'X-Admin-Token': tokenInput.value } })
  if (test.status === 200) {
    token.value = tokenInput.value
    localStorage.setItem('adminToken', tokenInput.value)
    tokenInput.value = ''
    emit('refresh')
  } else {
    tokenError.value = true
  }
}

function clearToken() {
  token.value = ''
  localStorage.removeItem('adminToken')
}

function resetCat() { catForm.value = { id: 0, name: '', sort_order: 0 } }
function resetNav() { navForm.value = { id: 0, title: '', url: '', category_id: 0, icon: '', description: '', sort_order: 0 } }

function editCat(c) { catForm.value = { ...c } }
function editNav(n) { navForm.value = { ...n } }

async function saveCategory() {
  const method = catForm.value.id ? 'PUT' : 'POST'
  const path = catForm.value.id ? `/categories/${catForm.value.id}` : '/categories'
  const res = await api(path, { method, headers: { 'Content-Type': 'application/json' }, body: JSON.stringify({ name: catForm.value.name, sort_order: catForm.value.sort_order }) })
  if (res.status === 401) { clearToken(); return }
  resetCat()
  emit('refresh')
}

async function saveNav() {
  const method = navForm.value.id ? 'PUT' : 'POST'
  const path = navForm.value.id ? `/navs/${navForm.value.id}` : '/navs'
  const res = await api(path, { method, headers: { 'Content-Type': 'application/json' }, body: JSON.stringify(navForm.value) })
  if (res.status === 401) { clearToken(); return }
  resetNav()
  emit('refresh')
}

async function delCat(id) {
  if (!confirm('确定删除该分类？其下导航也会被删除。')) return
  const res = await api(`/categories/${id}`, { method: 'DELETE' })
  if (res.status === 401) { clearToken(); return }
  emit('refresh')
}

async function delNav(id) {
  if (!confirm('确定删除该导航？')) return
  const res = await api(`/navs/${id}`, { method: 'DELETE' })
  if (res.status === 401) { clearToken(); return }
  emit('refresh')
}

function catName(cid) {
  const c = props.categories.find(x => x.id === cid)
  return c ? c.name : '-'
}

async function doExport() {
  const res = await api('/export', {})
  if (res.status === 401) { clearToken(); return }
  const data = await res.json()
  const blob = new Blob([JSON.stringify(data, null, 2)], { type: 'application/json' })
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = 'nav-backup.json'
  a.click()
  URL.revokeObjectURL(url)
}

async function doImport(e) {
  const file = e.target.files[0]
  if (!file) return
  const text = await file.text()
  const res = await api('/import', { method: 'POST', headers: { 'Content-Type': 'application/json' }, body: text })
  if (res.status === 401) { clearToken(); return }
  emit('refresh')
}

watch(() => props.modelValue, (v) => {
  if (v) {
    token.value = localStorage.getItem('adminToken') || ''
    emit('refresh')
  }
})
</script>
