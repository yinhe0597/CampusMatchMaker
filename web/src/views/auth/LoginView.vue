<template>
  <div style="max-width: 420px; margin: 80px auto;">
    <el-card shadow="never">
      <template #header>
        <h2 style="text-align: center; margin: 0; color: #303133;">登录</h2>
      </template>

      <el-form
        ref="formRef"
        :model="form"
        :rules="rules"
        label-width="0px"
        size="large"
        @keyup.enter="handleLogin"
        :loading="loading"
      >
        <el-form-item prop="student_id">
          <el-input
            v-model="form.student_id"
            placeholder="请输入学号"
            :prefix-icon="'User'"
          />
        </el-form-item>

        <el-form-item prop="password">
          <el-input
            v-model="form.password"
            type="password"
            placeholder="请输入密码"
            :prefix-icon="'Lock'"
            show-password
          />
        </el-form-item>

        <el-form-item>
          <el-button
            type="primary"
            style="width: 100%;"
            :loading="loading"
            @click="handleLogin"
          >
            {{ loading ? '登录中...' : '登录' }}
          </el-button>
        </el-form-item>

        <div style="text-align: center;">
          <router-link to="/register" style="color: #409eff; text-decoration: none;">
            还没有账号？去注册
          </router-link>
        </div>
      </el-form>
    </el-card>
  </div>
</template>

<script setup>
import { reactive, ref } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import { useUserStore } from '@/stores/user'

const router = useRouter()
const route = useRoute()
const userStore = useUserStore()

const formRef = ref(null)
const loading = ref(false)

const form = reactive({
  student_id: '',
  password: '',
})

const rules = {
  student_id: [
    { required: true, message: '请输入学号', trigger: 'blur' },
    { min: 6, max: 30, message: '学号长度应在 6-30 位', trigger: 'blur' },
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 8, message: '密码长度不能少于 8 位', trigger: 'blur' },
  ],
}

async function handleLogin() {
  if (!formRef.value) return

  const valid = await formRef.value.validate().catch(() => false)
  if (!valid) return

  loading.value = true
  try {
    await userStore.login({
      student_id: form.student_id,
      password: form.password,
    })
    ElMessage.success('登录成功')
    // 跳转到重定向页面或首页
    const redirect = route.query.redirect || '/polls'
    router.push(redirect)
  } catch (err) {
    ElMessage.error(err.message || '登录失败，请检查学号和密码')
  } finally {
    loading.value = false
  }
}
</script>
