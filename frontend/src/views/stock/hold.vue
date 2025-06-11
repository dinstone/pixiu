<template>
  <CommonPage>
    <div class="flex">
      <StockCard
        v-model:current-stock="currentStock"
        class="w-320 shrink-0"
        @refresh="refreshStock"
      />

      <HoldingCard :current-stock="currentStock" class="ml-40 w-0 flex-1" />
    </div>
  </CommonPage>
</template>

<script setup>
import HoldingCard from './components/HoldingCard.vue'
import StockCard from './components/StockCard.vue'

const currentStock = ref(null)
function refreshStock(stock) {
  currentStock.value = stock
}
</script>
