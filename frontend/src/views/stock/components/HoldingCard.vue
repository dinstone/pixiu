<template>
  <n-flex vertical>
    <template v-if="currentStock">
      <div class="flex justify-between">
        <h3 class="mb-12">
          持仓情况
        </h3>
      </div>
      <n-descriptions label-placement="left" bordered :column="3">
        <n-descriptions-item label="股票" :span="2">
          {{ currentStock.market }}  {{ currentStock.name }} ({{ currentStock.code }})
        </n-descriptions-item>
        <n-descriptions-item label="币种" :span="1">
          {{ currentStock.currency }}
        </n-descriptions-item>
        <n-descriptions-item label="持仓盈亏">
          {{ holding.profitLoss }}
        </n-descriptions-item>
        <n-descriptions-item label="税费合计">
          {{ holding.totalTaxFee }}
        </n-descriptions-item>
        <n-descriptions-item label="持股天数">
          {{ holding.holdingDays }}
        </n-descriptions-item>
        <n-descriptions-item label="持仓金额">
          {{ holding.amount }}
        </n-descriptions-item>
        <n-descriptions-item label="买入成本">
          {{ holding.costPrice }}
        </n-descriptions-item>
        <n-descriptions-item label="持仓数量">
          {{ holding.quantity }}
        </n-descriptions-item>
      </n-descriptions>

      <div class="mt-32 flex justify-between">
        <h3 class="mb-12">
          交易记录
        </h3>
        <NButton size="small" type="primary" @click="handleAddBtn">
          <i class="i-fe:plus mr-4 text-14" />
          新增
        </NButton>
      </div>

      <n-data-table
        :remote="true"
        :loading="loading"
        :columns="tradeColumns"
        :data="tradeList"
        :pagination="false"
        :bordered="false"
      />
    </template>
    <n-empty v-else class="h-450 f-c-c" size="large" description="请选择股票查看持仓详情" />
    <TradeEdit ref="modalRef" @refresh="initHolding" />
  </n-flex>
</template>

<script setup>
import { NButton } from 'naive-ui'
import api from '../api.js'
import TradeEdit from './TradeEdit.vue'

const props = defineProps({
  currentStock: {
    type: Object,
    default: () => null,
  },
})

// const currentStock = ref(null)

const holding = ref({})
const tradeList = ref([])
const loading = ref(false)

const tradeColumns = [
  { title: '类型', key: 'action', render(row) {
    if (row.action === 1) {
      return h('span', { style: 'color: red' }, '买入')
    }
    else {
      return h('span', { style: 'color: blue' }, '卖出')
    }
  } },
  { title: '数量', key: 'quantity' },
  { title: '价格', key: 'price' },
  { title: '金额', key: 'amount' },
  { title: '费用', key: 'taxFee' },
  {
    title: '时间',
    key: 'finishTime',
  },
  {
    title: '操作',
    key: 'actions',
    width: 200,
    align: 'center',
    fixed: 'right',
    render(row) {
      return [
        h(
          NButton,
          {
            size: 'small',
            type: 'primary',
            style: 'margin-left: 12px;',
            onClick: () => handleEditBtn(row),
          },
          {
            default: () => '编辑',
            icon: () => h('i', { class: 'i-material-symbols:edit-outline text-14' }),
          },
        ),

        h(
          NButton,
          {
            size: 'small',
            type: 'error',
            style: 'margin-left: 12px;',
            onClick: () => handleDeleteBtn(row.id),
          },
          {
            default: () => '删除',
            icon: () => h('i', { class: 'i-material-symbols:delete-outline text-14' }),
          },
        ),
      ]
    },
  },
]
async function initHolding() {
  loading.value = true
  try {
    // 查找持仓信息
    const res = await api.getHolding(props.currentStock.code)
    if (res?.data) {
      holding.value = { ...res?.data }
      const tres = await api.getTrades(res.data.id)
      tradeList.value = tres?.data || []
    }
    else {
      holding.value = {}
      tradeList.value = []
    }
    loading.value = false
  }
  catch (error) {
    console.error(error)
    loading.value = false
  }
}

watch(
  () => props.currentStock,
  async (v) => {
    await nextTick()

    if (v) {
      initHolding()
    }
  },
)

const modalRef = ref(null)
function handleAddBtn() {
  modalRef.value?.handleOpen({
    action: 'add',
    title: `新增交易 ${props.currentStock.name}`,
    row: { stockCode: props.currentStock.code },
    okText: '保存',
  })
}

function handleEditBtn(row) {
  modalRef.value?.handleOpen({
    action: 'edit',
    title: `编辑交易 ${props.currentStock.name}`,
    row: { stockCode: props.currentStock.code, ...row },
    okText: '保存',
  })
}

function handleDeleteBtn(id) {
  const d = $dialog.warning({
    content: '确定删除交易？',
    title: '提示',
    positiveText: '确定',
    negativeText: '取消',
    async onPositiveClick() {
      try {
        d.loading = true
        await api.deleteTrade(id)
        $message.success('删除成功')
        d.loading = false

        await initHolding()
      }
      catch (error) {
        console.error(error)
        d.loading = false
      }
    },
  })
}
</script>
