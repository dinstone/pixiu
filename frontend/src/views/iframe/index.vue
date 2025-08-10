<template>
  <AppPage full>
    <iframe :src="eurl" frameborder="0" class="wh-full" />
  </AppPage>
</template>

<script setup>
const route = useRoute()
const eurl = ref('')
if (route.query.site) {
  eurl.value = route.query.site
}
else if (route.meta.originPath) {
  eurl.value = route.meta.originPath
}
</script>
