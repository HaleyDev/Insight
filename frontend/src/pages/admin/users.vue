<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue'
import type { UserInfo } from '@/stores/auth'
import { createUser, deleteUser, listUsers, updateUser } from '@/utils/admin'
import { useAuthStore } from '@/stores/auth'

const auth = useAuthStore()

const users = ref<UserInfo[]>([])
const loading = ref(false)
const errorMsg = ref('')

// 创建对话框
const dialog = ref(false)
const submitting = ref(false)
const formError = ref('')
const form = reactive({
  username: '',
  email: '',
  password: '',
  confirmPassword: '',
})

// 编辑对话框
const editDialog = ref(false)
const editing = ref(false)
const editError = ref('')
const editForm = reactive({
  id: 0 as number,
  username: '',
  email: '',
  password: '',
  role: 'user',
})

// 删除确认
const deleteDialog = ref(false)
const deleting = ref(false)
const deleteError = ref('')
const deleteTarget = ref<UserInfo | null>(null)

const roleOptions = [
  { title: 'admin', value: 'admin' },
  { title: 'user', value: 'user' },
]

async function fetchList() {
  loading.value = true
  errorMsg.value = ''
  try {
    const data = await listUsers()
    users.value = Array.isArray(data) ? data : []
  }
  catch (e: any) {
    errorMsg.value = e?.message || '加载用户失败'
  }
  finally {
    loading.value = false
  }
}

function openCreate() {
  form.username = ''
  form.email = ''
  form.password = ''
  form.confirmPassword = ''
  formError.value = ''
  dialog.value = true
}

async function onSubmit() {
  formError.value = ''
  if (!form.username || !form.email || !form.password) {
    formError.value = '请填写完整信息'
    return
  }
  if (form.password !== form.confirmPassword) {
    formError.value = '两次密码不一致'
    return
  }
  submitting.value = true
  try {
    await createUser({ ...form })
    dialog.value = false
    await fetchList()
  }
  catch (e: any) {
    formError.value = e?.message || '创建失败'
  }
  finally {
    submitting.value = false
  }
}

function openEdit(item: UserInfo) {
  editForm.id = Number(item.id || 0)
  editForm.username = item.username || ''
  editForm.email = item.email || ''
  editForm.password = ''
  editForm.role = item.role || 'user'
  editError.value = ''
  editDialog.value = true
}

async function onEditSubmit() {
  editError.value = ''
  if (!editForm.id) {
    editError.value = '无效的用户 id'
    return
  }
  if (!editForm.username || !editForm.email) {
    editError.value = '用户名和邮箱不能为空'
    return
  }
  editing.value = true
  try {
    await updateUser(editForm.id, {
      username: editForm.username,
      email: editForm.email,
      role: editForm.role,
      // 密码留空表示不修改
      password: editForm.password || undefined,
    })
    editDialog.value = false
    await fetchList()
  }
  catch (e: any) {
    editError.value = e?.message || '更新失败'
  }
  finally {
    editing.value = false
  }
}

function openDelete(item: UserInfo) {
  deleteTarget.value = item
  deleteError.value = ''
  deleteDialog.value = true
}

async function onDeleteConfirm() {
  if (!deleteTarget.value?.id)
    return
  deleting.value = true
  deleteError.value = ''
  try {
    await deleteUser(Number(deleteTarget.value.id))
    deleteDialog.value = false
    deleteTarget.value = null
    await fetchList()
  }
  catch (e: any) {
    deleteError.value = e?.message || '删除失败'
  }
  finally {
    deleting.value = false
  }
}

const headers = [
  { title: 'ID', key: 'id', width: 80 },
  { title: '用户名', key: 'username' },
  { title: '邮箱', key: 'email' },
  { title: '角色', key: 'role', width: 120 },
  { title: '头像', key: 'avatar', sortable: false },
  { title: '操作', key: 'actions', sortable: false, width: 140, align: 'end' as const },
]

onMounted(fetchList)
</script>

<template>
  <VCard>
    <VCardItem>
      <template #title>
        用户管理
      </template>
      <template #subtitle>
        仅管理员可访问；账号仅能在此处由管理员创建。
      </template>
      <template #append>
        <VBtn
          color="primary"
          prepend-icon="ri-user-add-line"
          @click="openCreate"
        >
          添加用户
        </VBtn>
      </template>
    </VCardItem>

    <VAlert
      v-if="errorMsg"
      type="error"
      variant="tonal"
      density="compact"
      class="ma-4"
    >
      {{ errorMsg }}
    </VAlert>

    <VDataTable
      :headers="headers"
      :items="users"
      :loading="loading"
      class="text-no-wrap"
      item-value="id"
    >
      <template #item.role="{ item }">
        <VChip
          :color="item.role === 'admin' ? 'primary' : 'default'"
          size="small"
          label
        >
          {{ item.role || 'user' }}
        </VChip>
      </template>
      <template #item.avatar="{ item }">
        <VAvatar size="32">
          <VImg
            v-if="item.avatar"
            :src="item.avatar"
          />
          <span v-else>{{ (item.username || '?').charAt(0).toUpperCase() }}</span>
        </VAvatar>
      </template>
      <template #item.actions="{ item }">
        <VBtn
          icon
          variant="text"
          size="small"
          color="primary"
          @click="openEdit(item)"
        >
          <VIcon icon="ri-pencil-line" />
          <VTooltip activator="parent" location="top">
            编辑
          </VTooltip>
        </VBtn>
        <VBtn
          icon
          variant="text"
          size="small"
          color="error"
          :disabled="Number(item.id) === Number(auth.user?.id)"
          @click="openDelete(item)"
        >
          <VIcon icon="ri-delete-bin-line" />
          <VTooltip activator="parent" location="top">
            {{ Number(item.id) === Number(auth.user?.id) ? '不能删除当前登录账号' : '删除' }}
          </VTooltip>
        </VBtn>
      </template>
    </VDataTable>
  </VCard>

  <!-- 添加用户对话框 -->
  <VDialog
    v-model="dialog"
    max-width="480"
    persistent
  >
    <VCard>
      <VCardItem>
        <template #title>
          添加用户
        </template>
      </VCardItem>
      <VCardText>
        <VForm @submit.prevent="onSubmit">
          <VRow>
            <VCol cols="12">
              <VTextField
                v-model="form.username"
                label="用户名"
                autocomplete="off"
              />
            </VCol>
            <VCol cols="12">
              <VTextField
                v-model="form.email"
                label="邮箱"
                type="email"
                autocomplete="off"
              />
            </VCol>
            <VCol cols="12">
              <VTextField
                v-model="form.password"
                label="密码"
                type="password"
                autocomplete="new-password"
              />
            </VCol>
            <VCol cols="12">
              <VTextField
                v-model="form.confirmPassword"
                label="确认密码"
                type="password"
                autocomplete="new-password"
              />
            </VCol>
            <VCol
              v-if="formError"
              cols="12"
            >
              <VAlert
                type="error"
                variant="tonal"
                density="compact"
              >
                {{ formError }}
              </VAlert>
            </VCol>
          </VRow>
        </VForm>
      </VCardText>
      <VCardActions class="px-4 pb-4">
        <VSpacer />
        <VBtn
          variant="tonal"
          :disabled="submitting"
          @click="dialog = false"
        >
          取消
        </VBtn>
        <VBtn
          color="primary"
          :loading="submitting"
          @click="onSubmit"
        >
          创建
        </VBtn>
      </VCardActions>
    </VCard>
  </VDialog>

  <!-- 编辑用户对话框 -->
  <VDialog
    v-model="editDialog"
    max-width="480"
    persistent
  >
    <VCard>
      <VCardItem>
        <template #title>
          编辑用户 #{{ editForm.id }}
        </template>
        <template #subtitle>
          头像仅由用户本人设置，此处不可修改
        </template>
      </VCardItem>
      <VCardText>
        <VForm @submit.prevent="onEditSubmit">
          <VRow>
            <VCol cols="12">
              <VTextField
                v-model="editForm.username"
                label="用户名"
                autocomplete="off"
              />
            </VCol>
            <VCol cols="12">
              <VTextField
                v-model="editForm.email"
                label="邮箱"
                type="email"
                autocomplete="off"
              />
            </VCol>
            <VCol cols="12">
              <VSelect
                v-model="editForm.role"
                :items="roleOptions"
                label="角色"
              />
            </VCol>
            <VCol cols="12">
              <VTextField
                v-model="editForm.password"
                label="新密码（留空则不修改）"
                type="password"
                autocomplete="new-password"
              />
            </VCol>
            <VCol
              v-if="editError"
              cols="12"
            >
              <VAlert
                type="error"
                variant="tonal"
                density="compact"
              >
                {{ editError }}
              </VAlert>
            </VCol>
          </VRow>
        </VForm>
      </VCardText>
      <VCardActions class="px-4 pb-4">
        <VSpacer />
        <VBtn
          variant="tonal"
          :disabled="editing"
          @click="editDialog = false"
        >
          取消
        </VBtn>
        <VBtn
          color="primary"
          :loading="editing"
          @click="onEditSubmit"
        >
          保存
        </VBtn>
      </VCardActions>
    </VCard>
  </VDialog>

  <!-- 删除确认对话框 -->
  <VDialog
    v-model="deleteDialog"
    max-width="420"
    persistent
  >
    <VCard>
      <VCardItem>
        <template #title>
          删除用户
        </template>
      </VCardItem>
      <VCardText>
        <p class="mb-2">
          确认要删除用户
          <strong>{{ deleteTarget?.username }}</strong>
          (#{{ deleteTarget?.id }}) 吗？此操作不可恢复。
        </p>
        <VAlert
          v-if="deleteError"
          type="error"
          variant="tonal"
          density="compact"
        >
          {{ deleteError }}
        </VAlert>
      </VCardText>
      <VCardActions class="px-4 pb-4">
        <VSpacer />
        <VBtn
          variant="tonal"
          :disabled="deleting"
          @click="deleteDialog = false"
        >
          取消
        </VBtn>
        <VBtn
          color="error"
          :loading="deleting"
          @click="onDeleteConfirm"
        >
          删除
        </VBtn>
      </VCardActions>
    </VCard>
  </VDialog>
</template>
