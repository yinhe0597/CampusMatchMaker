<template>
  <div>
    <!-- 页面标题 -->
    <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 20px;">
      <h2 style="margin: 0;">我的班级</h2>
      <div style="display: flex; gap: 8px;">
        <el-button @click="showJoin = true" :icon="'Plus'">加入班级</el-button>
        <el-button type="primary" @click="showCreate = true" :icon="'Plus'">创建班级</el-button>
      </div>
    </div>

    <!-- 加载状态 -->
    <div v-if="classStore.loading && classStore.myClasses.length === 0" style="text-align: center; padding: 60px 0;">
      <el-skeleton :rows="3" animated />
    </div>

    <!-- 班级列表 -->
    <div v-else-if="classStore.myClasses.length > 0" style="display: grid; grid-template-columns: repeat(auto-fill, minmax(320px, 1fr)); gap: 16px;">
      <el-card
        v-for="cls in classStore.myClasses"
        :key="cls.id"
        shadow="hover"
        style="cursor: pointer;"
        @click="goToDetail(cls.id)"
      >
        <div style="display: flex; justify-content: space-between; align-items: flex-start;">
          <div>
            <h3 style="margin: 0 0 4px 0;">{{ cls.name }}</h3>
            <p style="margin: 0; color: #909399; font-size: 13px;">
              {{ cls.grade }}
              <template v-if="cls.department"> · {{ cls.department }}</template>
            </p>
          </div>
          <el-tag :type="roleTagType(cls.role)" size="small">
            {{ roleLabel(cls.role) }}
          </el-tag>
        </div>
        <el-divider style="margin: 12px 0;" />
        <div style="display: flex; justify-content: space-between; color: #909399; font-size: 13px;">
          <span>
            <el-icon><User /></el-icon>
            {{ cls.member_count }} 人
          </span>
          <span>
            <el-icon :color="cls.timetable_status === 1 ? '#67c23a' : '#e6a23c'">
              <Calendar />
            </el-icon>
            {{ cls.timetable_status === 1 ? '已录入课表' : '无课表' }}
          </span>
          <span>{{ cls.created_at?.slice(0, 10) }}</span>
        </div>
      </el-card>
    </div>

    <!-- 空状态 -->
    <el-empty v-else description="暂无班级，请创建或加入一个班级">
      <el-button type="primary" @click="showCreate = true">创建班级</el-button>
    </el-empty>

    <!-- 创建班级对话框 -->
    <el-dialog v-model="showCreate" title="创建班级" width="420px" :close-on-click-modal="false">
      <el-form ref="createFormRef" :model="createForm" :rules="createRules" label-width="80px">
        <el-form-item label="班级名称" prop="name">
          <el-input v-model="createForm.name" placeholder="如：计科2401班" maxlength="50" />
        </el-form-item>
        <el-form-item label="年级" prop="grade">
          <el-select v-model="createForm.grade" placeholder="选择年级" style="width: 100%;">
            <el-option v-for="y in gradeOptions" :key="y" :label="y" :value="y" />
          </el-select>
        </el-form-item>
        <el-form-item label="院系" prop="department">
          <el-input v-model="createForm.department" placeholder="如：计算机科学与技术学院（选填）" maxlength="50" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showCreate = false">取消</el-button>
        <el-button type="primary" :loading="creating" @click="handleCreate">创建</el-button>
      </template>
    </el-dialog>

    <!-- 加入班级对话框 -->
    <el-dialog v-model="showJoin" title="加入班级" width="400px" :close-on-click-modal="false">
      <template v-if="!joinStep">
        <el-form ref="joinFormRef" :model="joinForm" :rules="joinRules" label-width="0px">
          <el-form-item prop="invite_code">
            <el-input
              v-model="joinForm.invite_code"
              placeholder="请输入6位邀请码"
              maxlength="10"
              size="large"
              @keyup.enter="handleLookup"
            />
          </el-form-item>
        </el-form>
        <div style="text-align: center; margin-top: 12px;">
          <el-button type="primary" :loading="lookingUp" @click="handleLookup" style="width: 100%;">
            查找班级
          </el-button>
        </div>
      </template>
      <template v-else>
        <el-result icon="success" :title="'即将加入：' + joinInfo.name" :sub-title="joinInfo.grade + ' · 共 ' + joinInfo.member_count + ' 人'">
        </el-result>
      </template>
      <template #footer>
        <template v-if="joinStep">
          <el-button @click="resetJoin">重新输入</el-button>
          <el-button type="primary" :loading="joining" @click="handleJoin">确认加入</el-button>
        </template>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { reactive, ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { useClassStore } from '@/stores/class'
import { useUserStore } from '@/stores/user'

const router = useRouter()
const classStore = useClassStore()
const userStore = useUserStore()

// 年级选项
const gradeOptions = ['2026', '2025', '2024', '2023', '2022', '2021', '2020']

// ===== 创建班级 =====
const showCreate = ref(false)
const creating = ref(false)
const createFormRef = ref(null)
const createForm = reactive({
  name: '',
  grade: '',
  department: '',
})
const createRules = {
  name: [{ required: true, message: '请输入班级名称', trigger: 'blur' }],
  grade: [{ required: true, message: '请选择年级', trigger: 'change' }],
}

async function handleCreate() {
  if (!createFormRef.value) return
  const valid = await createFormRef.value.validate().catch(() => false)
  if (!valid) return

  creating.value = true
  try {
    await classStore.createClass({
      school_id: userStore.userInfo?.school_id || 1,
      name: createForm.name,
      grade: createForm.grade,
      department: createForm.department || undefined,
    })
    ElMessage.success('班级创建成功')
    showCreate.value = false
    createForm.name = ''
    createForm.grade = ''
    createForm.department = ''
  } catch (err) {
    ElMessage.error(err.message)
  } finally {
    creating.value = false
  }
}

// ===== 加入班级 =====
const showJoin = ref(false)
const joinStep = ref(0) // 0:输入邀请码, 1:确认加入
const lookingUp = ref(false)
const joining = ref(false)
const joinFormRef = ref(null)
const joinForm = reactive({
  invite_code: '',
})
const joinInfo = ref({})
const joinRules = {
  invite_code: [
    { required: true, message: '请输入邀请码', trigger: 'blur' },
    { len: 6, message: '邀请码为6位', trigger: 'blur' },
  ],
}

async function handleLookup() {
  if (!joinFormRef.value) return
  const valid = await joinFormRef.value.validate().catch(() => false)
  if (!valid) return

  lookingUp.value = true
  try {
    joinInfo.value = await classStore.lookupByCode(joinForm.invite_code)
    joinStep.value = 1
  } catch (err) {
    ElMessage.error(err.message || '未找到该班级，请检查邀请码')
  } finally {
    lookingUp.value = false
  }
}

async function handleJoin() {
  joining.value = true
  try {
    await classStore.joinClass(joinForm.invite_code)
    ElMessage.success('加入成功')
    showJoin.value = false
    resetJoin()
  } catch (err) {
    ElMessage.error(err.message)
  } finally {
    joining.value = false
  }
}

function resetJoin() {
  joinStep.value = 0
  joinForm.invite_code = ''
  joinInfo.value = {}
}

// ===== 工具函数 =====
function roleTagType(role) {
  if (role === 'owner') return 'danger'
  if (role === 'admin') return 'warning'
  return 'info'
}

function roleLabel(role) {
  if (role === 'owner') return '创建者'
  if (role === 'admin') return '管理员'
  return '成员'
}

function goToDetail(id) {
  router.push({ name: 'ClassDetail', params: { id } })
}

// 页面挂载时获取班级列表
onMounted(() => {
  classStore.fetchMyClasses()
})
</script>
