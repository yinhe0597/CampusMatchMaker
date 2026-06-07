<template>
  <div>
    <!-- 页面标题 -->
    <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 16px;">
      <div>
        <el-button text :icon="'ArrowLeft'" @click="goBack" style="margin-right: 8px;">返回详情</el-button>
        <h2 style="display: inline; margin: 0;">个人课表</h2>
      </div>
      <div v-if="classId" style="display: flex; gap: 8px;">
        <el-button v-if="canEdit" @click="handleEditClassTimetable">
          {{ hasClassTimetable ? '编辑课表' : '录入课表' }}
        </el-button>
        <el-button type="primary" @click="showAddPersonal = true">添加个人课程</el-button>
      </div>
    </div>

    <!-- 加载状态 -->
    <div v-if="loading" style="text-align: center; padding: 60px 0;">
      <el-skeleton :rows="8" animated />
    </div>

    <template v-else-if="personalTimetable?.entries?.length > 0">
      <!-- 统计信息 -->
      <div style="margin-bottom: 12px; font-size: 13px; color: #909399;">
        <el-tag size="small" type="info" style="margin-right: 8px;">继承 {{ personalTimetable.inherited_count }}</el-tag>
        <el-tag size="small" type="success">个人 {{ personalTimetable.personal_count }}</el-tag>
      </div>

      <!-- 课表网格 -->
      <div class="timetable-grid">
        <!-- 表头 -->
        <div class="tg-header tg-corner">节次</div>
        <div v-for="day in dayLabels" :key="'h-'+day" class="tg-header">{{ day }}</div>

        <!-- 行 -->
        <template v-for="period in 12" :key="'p-'+period">
          <div class="tg-period">
            <div class="tp-num">第{{ period }}节</div>
            <div class="tp-time">{{ periodTimes[period-1].start }}-{{ periodTimes[period-1].end }}</div>
          </div>
          <div
            v-for="dayIdx in 7"
            :key="'c-'+dayIdx+'-'+period"
            class="tg-cell"
            :class="getCellClass(gridData[dayIdx-1]?.[period-1])"
            :style="getCellStyle(gridData[dayIdx-1]?.[period-1])"
            @click="handleCellClick(gridData[dayIdx-1]?.[period-1])"
          >
            <template v-if="gridData[dayIdx-1]?.[period-1]?.entry">
              <div class="tc-course">{{ gridData[dayIdx-1][period-1].entry.course_name }}</div>
              <div v-if="gridData[dayIdx-1][period-1].entry.teacher" class="tc-teacher">
                {{ gridData[dayIdx-1][period-1].entry.teacher }}
              </div>
              <div v-if="gridData[dayIdx-1][period-1].entry.room" class="tc-room">
                {{ gridData[dayIdx-1][period-1].entry.room }}
              </div>
              <div v-if="gridData[dayIdx-1][period-1].entry.source" class="tc-source">
                {{ gridData[dayIdx-1][period-1].entry.source === 'inherited' ? '继承' : '个人' }}
              </div>
            </template>
          </div>
        </template>
      </div>
    </template>

    <!-- 空状态 -->
    <el-empty v-else description="暂无课表数据">
      <p style="color: #909399; font-size: 13px;">请先录入班级公共课表，或添加个人课程</p>
      <el-button v-if="canEdit && classId" type="primary" @click="handleEditClassTimetable">录入课表</el-button>
    </el-empty>

    <!-- 添加个人课程对话框 -->
    <el-dialog v-model="showAddPersonal" title="添加个人课程" width="420px" :close-on-click-modal="false">
      <el-form ref="addFormRef" :model="addForm" :rules="addRules" label-width="90px">
        <el-form-item label="星期" prop="day_of_week">
          <el-select v-model="addForm.day_of_week" placeholder="选择星期" style="width: 100%;">
            <el-option v-for="(label, idx) in dayLabels" :key="idx+1" :label="label" :value="idx+1" />
          </el-select>
        </el-form-item>
        <el-form-item label="开始节次" prop="period_start">
          <el-select v-model="addForm.period_start" placeholder="选择节次" style="width: 100%;">
            <el-option v-for="p in 12" :key="p" :label="'第'+p+'节 ('+periodTimes[p-1].start+')'" :value="p" />
          </el-select>
        </el-form-item>
        <el-form-item label="结束节次" prop="period_end">
          <el-select v-model="addForm.period_end" placeholder="选择节次" style="width: 100%;">
            <el-option v-for="p in 12" :key="p" :label="'第'+p+'节 ('+periodTimes[p-1].end+')'" :value="p" />
          </el-select>
        </el-form-item>
        <el-form-item label="课程名称" prop="course_name">
          <el-input v-model="addForm.course_name" placeholder="如：英语选修" maxlength="50" />
        </el-form-item>
        <el-form-item label="教师">
          <el-input v-model="addForm.teacher" placeholder="（选填）" maxlength="20" />
        </el-form-item>
        <el-form-item label="教室">
          <el-input v-model="addForm.room" placeholder="（选填）" maxlength="20" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showAddPersonal = false">取消</el-button>
        <el-button type="primary" :loading="adding" @click="handleAddPersonal">添加</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { reactive, ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { useTimetableStore } from '@/stores/timetable'
import { useUserStore } from '@/stores/user'
import { PERIOD_TIMES, DAY_LABELS, getPeriodLabel, getPeriodRowSpan } from '@/utils/time'

const route = useRoute()
const router = useRouter()
const timetableStore = useTimetableStore()
const userStore = useUserStore()

const classId = computed(() => Number(route.query.class_id) || 0)
const loading = ref(false)
const adding = ref(false)
const showAddPersonal = ref(false)

const personalTimetable = computed(() => timetableStore.personalTimetable)
const dayLabels = DAY_LABELS
const periodTimes = PERIOD_TIMES

// 权限: owner/admin 可编辑班级课表
const canEdit = computed(() => {
  const role = route.query.role
  return role === 'owner' || role === 'admin'
})

const hasClassTimetable = computed(() => {
  return personalTimetable.value?.inherited_count > 0
})

// 构建 7x12 网格
const gridData = computed(() => {
  const grid = Array.from({ length: 7 }, () => Array(12).fill(null))
  const entries = personalTimetable.value?.entries || []

  for (const entry of entries) {
    const dayIdx = entry.day_of_week - 1
    const periodIdx = entry.period_start - 1
    const rowspan = entry.period_end - entry.period_start + 1

    // 填充起始单元格
    grid[dayIdx][periodIdx] = { entry, rowspan }

    // 标记后续被占用的行
    for (let r = 1; r < rowspan; r++) {
      if (periodIdx + r < 12) {
        grid[dayIdx][periodIdx + r] = { hidden: true }
      }
    }
  }

  return grid
})

function getCellClass(cell) {
  if (!cell || cell.hidden) return 'tg-hidden'
  const source = cell.entry?.source
  return {
    'tg-inherited': source === 'inherited',
    'tg-personal': source === 'personal',
  }
}

function getCellStyle(cell) {
  if (!cell || !cell.rowspan || cell.rowspan <= 1) return {}
  return { gridRowEnd: `span ${cell.rowspan}` }
}

// 添加个人课程表单
const addFormRef = ref(null)
const addForm = reactive({
  day_of_week: '',
  period_start: '',
  period_end: '',
  course_name: '',
  teacher: '',
  room: '',
})
const addRules = {
  day_of_week: [{ required: true, message: '请选择星期', trigger: 'change' }],
  period_start: [{ required: true, message: '请选择开始节次', trigger: 'change' }],
  period_end: [{ required: true, message: '请选择结束节次', trigger: 'change' }],
  course_name: [{ required: true, message: '请输入课程名称', trigger: 'blur' }],
}

async function handleAddPersonal() {
  if (!addFormRef.value) return
  const valid = await addFormRef.value.validate().catch(() => false)
  if (!valid) return

  if (addForm.period_start > addForm.period_end) {
    ElMessage.warning('开始节次不能大于结束节次')
    return
  }

  adding.value = true
  try {
    await timetableStore.createPersonalTimetable({
      class_id: classId.value,
      day_of_week: Number(addForm.day_of_week),
      period_start: Number(addForm.period_start),
      period_end: Number(addForm.period_end),
      course_name: addForm.course_name,
      teacher: addForm.teacher || undefined,
      room: addForm.room || undefined,
    })
    ElMessage.success('添加成功')
    showAddPersonal.value = false
    addForm.day_of_week = ''
    addForm.period_start = ''
    addForm.period_end = ''
    addForm.course_name = ''
    addForm.teacher = ''
    addForm.room = ''
  } catch (err) {
    ElMessage.error(err.message)
  } finally {
    adding.value = false
  }
}

function handleEditClassTimetable() {
  router.push({
    name: 'TimetableEdit',
    query: { class_id: classId.value },
  })
}

async function handleCellClick(cell) {
  if (!cell || !cell.entry || cell.entry.source !== 'personal') return
  try {
    await ElMessageBox.confirm(
      `确定要删除「${cell.entry.course_name}」吗？`,
      '删除个人课程',
      { confirmButtonText: '删除', cancelButtonText: '取消', type: 'warning' }
    )
    await timetableStore.deletePersonalTimetable(cell.entry.id, classId.value)
    ElMessage.success('已删除')
  } catch {
    // 取消操作
  }
}

function goBack() {
  router.push({ name: 'ClassDetail', params: { id: classId.value } })
}

// 初始化
onMounted(async () => {
  if (classId.value) {
    loading.value = true
    try {
      await timetableStore.fetchPersonalTimetable(classId.value)
    } finally {
      loading.value = false
    }
  }
})
</script>

<style scoped>
.timetable-grid {
  display: grid;
  grid-template-columns: 80px repeat(7, 1fr);
  gap: 1px;
  background: #e4e7ed;
  border: 1px solid #e4e7ed;
  border-radius: 6px;
  overflow: hidden;
}

.tg-header {
  background: #f5f7fa;
  padding: 10px 6px;
  text-align: center;
  font-weight: 600;
  font-size: 13px;
  color: #303133;
}

.tg-corner {
  background: #ebeef5;
}

.tg-period {
  background: #f5f7fa;
  padding: 6px 4px;
  text-align: center;
  font-size: 11px;
  color: #606266;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
}

.tp-num {
  font-weight: 600;
  font-size: 12px;
}

.tp-time {
  font-size: 10px;
  color: #909399;
  margin-top: 2px;
}

.tg-cell {
  background: #fff;
  min-height: 48px;
  padding: 4px 6px;
  font-size: 12px;
  display: flex;
  flex-direction: column;
  justify-content: center;
}

.tg-hidden {
  opacity: 0;
  min-height: 0;
  padding: 0;
}

.tg-inherited {
  background: #ecf5ff;
  border-left: 3px solid #409eff;
}

.tg-personal {
  background: #f0f9eb;
  border-left: 3px solid #67c23a;
}

.tc-course {
  font-weight: 600;
  color: #303133;
  font-size: 13px;
  line-height: 1.3;
}

.tc-teacher,
.tc-room {
  color: #909399;
  font-size: 11px;
  margin-top: 2px;
}

.tc-source {
  position: absolute;
  top: 2px;
  right: 4px;
  font-size: 10px;
  color: #c0c4cc;
}
</style>
