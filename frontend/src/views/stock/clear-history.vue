<template>
  <CommonPage>
    <n-flex vertical>
      <n-descriptions label-placement="left" bordered :column="3">
        <n-descriptions-item label="股票" :span="2">
          {{ stockClear.market }}  {{ stockClear.name }} ({{ stockClear.code }})
        </n-descriptions-item>
        <n-descriptions-item label="币种" :span="1">
          {{ stockClear.currency }}
        </n-descriptions-item>
        <n-descriptions-item label="统计范围">
          {{ stockClear.startTime }} - {{ stockClear.finishTime }}
        </n-descriptions-item>
        <n-descriptions-item label="清仓次数">
          {{ stockClear.totalCount }}
        </n-descriptions-item>
        <n-descriptions-item label="盈亏金额">
          <span :style="{ color: stockClear.profitLoss > 0 ? 'red' : 'blue' }">
            {{ stockClear.profitLoss }}
          </span>
        </n-descriptions-item>
      </n-descriptions>

      <n-card title="清仓记录">
        <n-grid x-gap="12" :cols="2">
          <n-gi>
            <n-data-table
              :remote="true"
              :columns="investColumns"
              :data="investList"
              :pagination="false"
              :bordered="true"
              :row-props="rowProps"
            />
          </n-gi>
          <n-gi>
            <n-data-table
              :remote="true"
              :loading="loading"
              :columns="tradeColumns"
              :data="tradeList"
              :pagination="false"
              :bordered="true"
            />
          </n-gi>
        </n-grid>
      </n-card>
    </n-flex>
  </CommonPage>
</template>

<script setup>
import { onMounted } from 'vue'
import api from './api.js'

const investColumns = [
  { title: '投资时间', render: row => `${row.openTime?.substring(0, 10)} - ${row.closeTime?.substring(0, 10)}` },
  { title: '持股天数', key: 'holdingDays' },
  { title: '税费合计', key: 'totalTaxFee', render: row => row.totalTaxFee.toFixed(2) },
  { title: '清仓盈亏', key: 'profitLoss', render(row) {
    if (row.profitLoss > 0) {
      return h('span', { style: 'color: red' }, row.profitLoss.toFixed(2))
    }
    else {
      return h('span', { style: 'color: blue' }, row.profitLoss.toFixed(2))
    }
  } },
]

const route = useRoute()
const { stockCode, startTime, finishTime } = route.query

const stockClear = ref({})
const investList = ref([])
const loading = ref(false)
async function initStockClear() {
  loading.value = true
  try {
    const res = await api.getStockClear(stockCode, startTime, finishTime)
    if (res?.data) {
      stockClear.value = { ...res.data?.stock, ...res.data?.stats } || {}
      investList.value = res?.data.invests || []
    }
    loading.value = false
  }
  catch (error) {
    console.error(error)
    loading.value = false
  }
}

function rowProps(row) {
  return {
    style: 'cursor: pointer;',
    onClick: () => {
      initTradeList(row.id)
    },
  }
}

onMounted(() => {
  initStockClear()
})

const tradeList = ref([])
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
  { title: '时间', key: 'finishTime' },
]
async function initTradeList(investId) {
  loading.value = true
  try {
    const tres = await api.getTrades(investId)
    tradeList.value = tres?.data || []
    loading.value = false
  }
  catch (error) {
    console.error(error)
    loading.value = false
  }
}
</script>
