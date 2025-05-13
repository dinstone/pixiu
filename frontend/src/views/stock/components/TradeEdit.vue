<template>
  <MeModal ref="modalRef">
    <n-form
      ref="formRef"
      label-placement="left"
      require-mark-placement="left"
      :label-width="100"
      :model="formModel"
    >
      <n-grid :cols="24" :x-gap="24">
        <n-form-item-gi :span="12" label="类型" path="action" :rule="numberRule">
          <n-select v-model:value="formModel.action" :options="typeOptions" />
        </n-form-item-gi>

        <n-form-item-gi :span="12" path="quantity" :rule="numberRule">
          <template #label>
            数量
          </template>
          <n-input-number v-model:value="formModel.quantity" />
        </n-form-item-gi>
        <n-form-item-gi :span="12" path="price" :rule="numberRule">
          <template #label>
            价格
          </template>
          <n-input-number v-model:value="formModel.price" class="n-input n-input--resizable n-input--stateful" />
        </n-form-item-gi>
        <n-form-item-gi :span="12" path="taxFee">
          <template #label>
            费用
          </template>
          <n-input-number v-model:value="formModel.taxFee" />
        </n-form-item-gi>
        <n-form-item-gi :span="12" label="时间" path="finishTime" :rule="stringRule">
          <n-date-picker
            v-model:formatted-value="formModel.finishTime" value-format="yyyy-MM-dd HH:mm:ss"
            type="datetime" clearable class="n-input n-input--resizable n-input--stateful"
          />
        </n-form-item-gi>
      </n-grid>
    </n-form>
  </MeModal>
</template>

<script setup>
import { MeModal } from '@/components'
import { useForm, useModal } from '@/composables'
import api from '../api'

const emit = defineEmits(['refresh'])

const typeOptions = computed(() => {
  return [{ label: '买入', value: 1 }, { label: '卖出', value: -1 }]
})

const numberRule = {
  required: true,
  type: 'number',
  message: '此为必填项',
  trigger: ['blur', 'change'],
}
const stringRule = {
  required: true,
  type: 'string',
  message: '此为必填项',
  trigger: ['blur', 'change'],
}

const [formRef, formModel, validation] = useForm()
const [modalRef, okLoading] = useModal()

const modalAction = ref('')
function handleOpen(options = {}) {
  const { action, row = {}, ...rest } = options
  modalAction.value = action
  formModel.value = { ...row }
  modalRef.value.open({ ...rest, onOk: onSave })
}

async function onSave() {
  await validation()

  okLoading.value = true
  try {
    let res
    if (modalAction.value === 'add') {
      res = await api.addTrade(formModel.value)
    }
    else if (modalAction.value === 'edit') {
      res = await api.saveTrade(formModel.value)
    }
    okLoading.value = false

    if (res?.code === 0) {
      $message.success('保存成功')
      emit('refresh')
      return true
    }
    else {
      $message.error(res?.mesg)
      return false
    }
  }
  catch (error) {
    console.error(error)
    okLoading.value = false
    return false
  }
}

defineExpose({
  handleOpen,
})
</script>
