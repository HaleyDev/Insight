<script setup lang="ts">
const route = useRoute()
const router = useRouter()
const from = computed(() => (route.query.from as string) || '')
</script>

<template>
  <div class="d-flex flex-column align-center justify-center text-center pa-8" style="min-height: 60vh;">
    <VIcon
      icon="ri-shield-cross-line"
      size="80"
      color="error"
      class="mb-4"
    />
    <h1 class="text-h3 mb-2">
      403
    </h1>
    <h2 class="text-h5 mb-2">
      无权限访问
    </h2>
    <p class="text-body-1 mb-6">
      当前账号没有权限访问该页面
      <span v-if="from" class="text-disabled">（{{ from }}）</span>
    </p>
    <div class="d-flex gap-3">
      <VBtn
        color="primary"
        @click="router.replace('/dashboard')"
      >
        返回首页
      </VBtn>
      <VBtn
        variant="outlined"
        @click="router.back()"
      >
        返回上一页
      </VBtn>
    </div>
  </div>
</template>
