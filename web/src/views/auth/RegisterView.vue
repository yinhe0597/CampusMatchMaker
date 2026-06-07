<template>
  <div style="max-width: 420px; margin: 80px auto;">
    <el-card shadow="never">
      <template #header>
        <h2 style="text-align: center; margin: 0; color: #303133;">注册</h2>
      </template>

      <el-form
        ref="formRef"
        :model="form"
        :rules="rules"
        label-width="0px"
        size="large"
        @keyup.enter="handleRegister"
      >
        <el-form-item prop="student_id">
          <el-input
            v-model="form.student_id"
            placeholder="请输入学号"
            :prefix-icon="'User'"
          />
        </el-form-item>

        <el-form-item prop="nickname">
          <el-input
            v-model="form.nickname"
            placeholder="请输入昵称"
            :prefix-icon="'Edit'"
          />
        </el-form-item>

        <el-form-item prop="school_id">
          <el-select
            v-model="form.school_id"
            placeholder="请选择学校"
            style="width: 100%;"
          >
            <el-option
              v-for="school in schools"
              :key="school.id"
              :label="school.name"
              :value="school.id"
            />
          </el-select>
        </el-form-item>

        <el-form-item prop="password">
          <el-input
            v-model="form.password"
            type="password"
            placeholder="请输入密码（至少8位）"
            :prefix-icon="'Lock'"
            show-password
          />
        </el-form-item>

        <el-form-item prop="confirmPassword">
          <el-input
            v-model="form.confirmPassword"
            type="password"
            placeholder="请确认密码"
            :prefix-icon="'Lock'"
            show-password
          />
        </el-form-item>

        <el-form-item>
          <el-button
            type="primary"
            style="width: 100%;"
            :loading="loading"
            @click="handleRegister"
          >
            {{ loading ? '注册中...' : '注册' }}
          </el-button>
        </el-form-item>

        <div style="text-align: center;">
          <router-link to="/login" style="color: #409eff; text-decoration: none;">
            已有账号？去登录
          </router-link>
        </div>
      </el-form>
    </el-card>
  </div>
</template>

<script setup>
import { reactive, ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { useUserStore } from '@/stores/user'
import request from '@/api/request'

const router = useRouter()
const userStore = useUserStore()

const formRef = ref(null)
const loading = ref(false)
const schools = ref([])

const form = reactive({
  student_id: '',
  nickname: '',
  school_id: null,
  password: '',
  confirmPassword: '',
})

const validatePass2 = (rule, value, callback) => {
  if (value !== form.password) {
    callback(new Error('两次输入的密码不一致'))
  } else {
    callback()
  }
}

const rules = {
  student_id: [
    { required: true, message: '请输入学号', trigger: 'blur' },
    { min: 6, max: 30, message: '学号长度应在 6-30 位', trigger: 'blur' },
  ],
  nickname: [
    { required: true, message: '请输入昵称', trigger: 'blur' },
    { min: 1, max: 50, message: '昵称长度应在 1-50 位', trigger: 'blur' },
  ],
  school_id: [
    { required: true, message: '请选择学校', trigger: 'change' },
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 8, max: 50, message: '密码长度应在 8-50 位', trigger: 'blur' },
  ],
  confirmPassword: [
    { required: true, message: '请确认密码', trigger: 'blur' },
    { validator: validatePass2, trigger: 'blur' },
  ],
}

// 获取学校列表
onMounted(async () => {
  // 暂时用硬编码学校数据，后续可从接口获取
  schools.value = [
    { id: 1, name: '测试大学' },
  ]
})

async function handleRegister() {
  if (!formRef.value) return

  const valid = await formRef.value.validate().catch(() => false)
  if (!valid) return

  loading.value = true
  try {
    await userStore.register({
      student_id: form.student_id,
      nickname: form.nickname,
      school_id: form.school_id,
      password: form.password,
    })
    ElMessage.success('注册成功')
    router.push('/polls')
  } catch (err) {
    ElMessage.error(err.message || '注册失败')
  } finally {
    loading.value = false
  }
}
</script>
