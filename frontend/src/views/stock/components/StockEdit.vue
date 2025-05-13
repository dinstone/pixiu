<template>
  <MeModal ref="modalRef">
    <n-form
      ref="modalFormRef"
      label-placement="left"
      require-mark-placement="left"
      :label-width="100"
      :model="modalForm"
    >
      <n-grid :cols="24" :x-gap="24">
        <n-form-item-gi :span="12" path="name" :rule="rules.required">
          <template #label>
            <QuestionLabel label="名称" content="股票全称" />
          </template>
          <n-input v-model:value="modalForm.name" />
        </n-form-item-gi>
        <n-form-item-gi :span="12" path="code" :rule="rules.required">
          <template #label>
            <QuestionLabel label="编码" content="股票编码不可修改" />
          </template>
          <n-input v-model:value="modalForm.code" />
        </n-form-item-gi>
        <n-form-item-gi :span="12" label="股市" path="market" :rule="rules.required">
          <n-select v-model:value="modalForm.market" :options="marketOptions" />
        </n-form-item-gi>
        <n-form-item-gi :span="12" label="币种" path="currency" :rule="rules.required">
          <n-select v-model:value="modalForm.currency" :options="currencyOptions" />
        </n-form-item-gi>
      </n-grid>
    </n-form>
  </MeModal>
</template>

<script setup>
import { MeModal } from '@/components'
import { useForm, useModal } from '@/composables'
import api from '../api'
import QuestionLabel from './QuestionLabel.vue'

const emit = defineEmits(['refresh'])

const marketOptions = computed(() => {
  return [{ label: 'A股', value: 'A股' }, { label: 'H股', value: 'H股' }]
})

const currencyOptions = [
  { label: '人民币', value: '人民币' },
  { label: '港币', value: '港币' },
  { label: '美元', value: '美元' },
]

const defaultForm = { enable: true, show: true, layout: '', amount: 0 }
const [modalFormRef, modalForm, validation, rules] = useForm()
const [modalRef, okLoading] = useModal()

const modalAction = ref('')
const parentIdDisabled = ref(false)
function handleOpen(options = {}) {
  const { action, row = {}, ...rest } = options
  modalAction.value = action
  modalForm.value = { ...defaultForm, ...row }
  parentIdDisabled.value = !!row.parentId && row.type === 'BUTTON'
  modalRef.value.open({ ...rest, onOk: onSave })
}

async function onSave() {
  await validation()
  okLoading.value = true
  try {
    let newFormData
    if (!modalForm.value.parentId)
      modalForm.value.parentId = null
    if (modalAction.value === 'add') {
      const res = await api.addStock(modalForm.value)
      newFormData = res.data
    }
    else if (modalAction.value === 'edit') {
      await api.saveStock(modalForm.value)
    }
    okLoading.value = false
    $message.success('保存成功')
    emit('refresh', modalAction.value === 'add' ? newFormData : modalForm.value)
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
