<template>
  <div class="contact">
    <h1>联系作者</h1>
    <el-card class="contact-card">
      <el-form :model="form" label-width="100px">
        <el-form-item label="姓名">
          <el-input v-model="form.name" />
        </el-form-item>
        <el-form-item label="邮箱">
          <el-input v-model="form.email" />
        </el-form-item>
        <el-form-item label="留言">
          <el-input
            v-model="form.message"
            type="textarea"
            :rows="4"
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" :loading="submitting" @click="submitForm">提交</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { ElMessage } from 'element-plus'
import request from '@/utils/request'
import type { ApiResponse } from '@/types/api'

const form = ref({
  name: '',
  email: '',
  message: ''
})

const submitting = ref(false)

const submitForm = async () => {
  if (!form.value.name || !form.value.message) {
    ElMessage.warning('请填写姓名和留言')
    return
  }
  try {
    submitting.value = true
    await request.post<ApiResponse<void>>('/contact-messages', {
      name: form.value.name,
      email: form.value.email,
      content: form.value.message
    })
    ElMessage.success('提交成功')
    form.value = { name: '', email: '', message: '' }
  } catch {
    // 错误提示已由 request 拦截器统一处理
  } finally {
    submitting.value = false
  }
}
</script>

<style scoped lang="scss">
@use "@/styles/variables" as *;

.contact {
  padding: 20px;
  max-width: $content-max-width;
  margin: 0 auto;
}
.contact-card {
  max-width: 600px;
  margin: 20px auto 0;
  border-color: $border-light;
  border-radius: $radius-xl;

  :deep(.el-card__body) {
    padding: $space-xl;
  }
}
</style>
