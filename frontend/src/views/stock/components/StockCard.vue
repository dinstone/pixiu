<template>
  <n-spin size="small" :show="treeLoading">
    <n-space vertical :size="12">
      <div class="flex">
        <n-input v-model:value="pattern" placeholder="搜索(股票编码、名称)" clearable />
        <NButton class="ml-12" type="primary" @click="handleAdd()">
          <i class="i-material-symbols:add mr-4 text-14" />
          新增
        </NButton>
      </div>

      <n-tree
        :show-irrelevant-nodes="false"
        :pattern="pattern"
        :data="treeData"
        :selected-keys="[currentStock?.code]"
        :render-prefix="renderPrefix"
        :render-label="renderLabel"
        :render-suffix="renderSuffix"
        :on-update:selected-keys="onSelect"
        :filter="searchFilter"
        key-field="code"
        label-field="name code"

        block-line default-expand-all
      />
    </n-space>

    <StockEdit ref="modalRef" @refresh="initStockTree" />
  </n-spin>
</template>

<script setup>
import { NButton } from 'naive-ui'
import { withModifiers } from 'vue'
import api from '../api'
import StockEdit from './StockEdit.vue'

defineProps({
  currentStock: {
    type: Object,
    default: () => null,
  },
})

const emit = defineEmits(['refresh', 'update:currentStock'])

const treeData = ref([])
const treeLoading = ref(false)
async function initStockTree() {
  treeLoading.value = true
  const res = await api.getStocks()
  treeData.value = res?.data || []
  treeLoading.value = false
}
initStockTree()

const pattern = ref('')
function searchFilter(pattern, node) {
  if (!pattern) {
    return true
  }
  const name = node.name || ''
  const code = node.code || ''
  return name.toLowerCase().includes(pattern.toLowerCase()) || code.toLowerCase().includes(pattern.toLowerCase())
}

const modalRef = ref(null)
async function handleAdd(data = {}) {
  modalRef.value?.handleOpen({
    action: 'add',
    title: '新增股票',
    row: { type: 'MENU', ...data },
    okText: '保存',
  })
}

function handleEdit(item = {}) {
  modalRef.value?.handleOpen({
    action: 'edit',
    title: `编辑股票 - ${item.name}`,
    row: item,
    okText: '保存',
  })
}

function handleDelete(item) {
  $dialog.confirm({
    content: `确认删除【${item.name}】？`,
    async confirm() {
      try {
        $message.loading('正在删除', { key: 'deleteMenu' })
        await api.deleteStock(item.code)
        $message.success('删除成功', { key: 'deleteMenu' })
        emit('refresh')
        emit('update:currentStock', null)
        initStockTree()
      }
      catch (error) {
        console.error(error)
        $message.destroy('deleteMenu')
      }
    },
  })
}

function onSelect(keys, option, { action, node }) {
  emit('update:currentStock', action === 'select' ? node : null)
}

function renderPrefix({ option }) {
  return option.market
}

function renderLabel({ option }) {
  const label = `${option.name}(${option.code})`
  return h('span', { class: 'text-primary' }, label)
}

function renderSuffix({ option }) {
  return [
    h(
      NButton,
      {
        text: true,
        type: 'primary',
        title: '编辑股票',
        size: 'tiny',
        onClick: withModifiers(() => handleEdit(option), ['stop']),
      },
      { default: () => '编辑' },
    ),

    h(
      NButton,
      {
        text: true,
        type: 'error',
        title: '删除股票',
        size: 'tiny',
        style: 'margin-left: 12px;',
        onClick: withModifiers(() => handleDelete(option), ['stop']),
      },
      { default: () => '删除' },
    ),
  ]
}
</script>
