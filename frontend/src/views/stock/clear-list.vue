<template>
  <CommonPage>
    <MeCrud
      ref="$table"
      v-model:query-items="queryItems"
      :scroll-x="1200"
      :columns="columns"
      :get-data="api.getClearList"
      :row-props="rowProps"
    >
      <MeQueryItem label="起始日期" :label-width="80">
        <n-date-picker
          v-model:formatted-value="queryItems.startTime" value-format="yyyy-MM-dd"
          type="date" clearable class="n-input n-input--resizable n-input--stateful"
        />
      </MeQueryItem>
      <MeQueryItem label="截止日期" :label-width="80">
        <n-date-picker
          v-model:formatted-value="queryItems.finishTime" value-format="yyyy-MM-dd"
          type="date" clearable class="n-input n-input--resizable n-input--stateful"
        />
      </MeQueryItem>
    </MeCrud>
  </CommonPage>
</template>

<script setup>
import { MeCrud, MeQueryItem } from '@/components'
import api from './api'

const $table = ref(null)
/** QueryBar筛选参数（可选） */
const queryItems = ref({})

const columns = [
  { title: '股票编码', key: 'stockCode' },
  { title: '股票名称', key: 'stockName' },
  { title: '投资盈亏', key: 'profitLoss' },
  { title: '投资次数', key: 'totalCount' },
  { title: '盈利次数', key: 'profitCount' },
  { title: '亏损次数', key: 'lossCount' },
]

onMounted(() => {
  $table.value?.handleSearch()
})

const router = useRouter()
function rowProps(row) {
  return {
    style: 'cursor: pointer;',
    onClick: () => {
      // $message.info(row.stockName)
      router.push({ path: '/stock/clear/history', query: { stockCode: row.stockCode, startTime: queryItems.value.startTime, finishTime: queryItems.value.finishTime } })
    },
  }
}
</script>
