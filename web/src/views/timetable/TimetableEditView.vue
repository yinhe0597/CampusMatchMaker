<template>
  <div>
    <el-button text :icon="'ArrowLeft'" @click="goBack" style="margin-bottom: 12px;">返回课表</el-button>

    <el-card shadow="never">
      <template #header>
        <div style="display: flex; justify-content: space-between; align-items: center;">
          <h3 style="margin: 0;">
            {{ isEdit ? '编辑班级课表' : '录入班级课表' }}
          </h3>
          <div>
            <el-button @click="addRow">添加一行</el-button>
            <el-button type="primary" :loading="submitting" @click="handleSubmit">
              {{ isEdit ? '保存修改' : '录入' }}
            </el-button>
          </div>
        </div>
      </template>

      <el-alert
        v-if="entries.length === 0"
        title="请添加班级公共课表条目，所有成员将自动继承"
        type="info"
        :closable="false"
        show-icon
        style="margin-bottom: 16px;"
      />

      <div v-for="(row, idx) in entries" :key="idx" style="margin-bottom: 12px;">
        <el-card shadow="hover" style="position: relative;">
          <el-button
            text
            type="danger"
            style="position: absolute; top: 8px; right: 8px;"
            @click="removeRow(idx)"
          >
            删除
          </el-button>
          <el-form :model="row" label-width="80px" size="small">
            <el-row :gutter="12">
              <el-col :span="6">
                <el-form-item label="星期">
                  <el-select v-model="row.day_of_week" placeholder="选择">
                    <el-option v-for="(label, i) in dayLabels" :key="i+1" :label="label" :value="i+1" />
                  </el-select>
                </el-form-item>
              </el-col>
              <el-col :span="4">
                <el-form-item label="开始">
                  <el-select v-model="row.period_start" placeholder="节次">
                    <el-option v-for="p in 12" :key="p" :label="'第'+p+'节'" :value="p" />
                  </el-select>
                </el-form-item>
              </el-col>
              <el-col :span="4">
                <el-form-item label="结束">
                  <el-select v-model="row.period_end" placeholder="节次">
                    <el-option v-for="p in 12" :key="p" :label="'第'+p+'节'" :value="p" />
                  </el-select>
                </el-form-item>
              </el-col>
              <el-col :span="5">
                <el-form-item label="课程">
                  <el-input v-model="row.course_name" placeholder="课程名称" />
                </el-form-item>
              </el-col>
              <el-col :span="3">
                <el-form-item label="教师">
                  <el-input v-model="row.teacher" placeholder="（选填）" />
                </el-form-item>
              </el-col>
              <el-col :span="2">
                <el-form-item label="教室">
                  <el-input v-model="row.room" placeholder="（选填）" />
                </el-form-item>
              </el-col>
            </el-row>
          </el-form>
        </el-card>
      </div>

      <el-empty v-if="entries.length === 0" description="暂无条目，点击上方「添加一行」开始录入" />
    </el-card>

    <!-- 前提条件检查 -->
    <el-dialog v-model="showPrereq" title="无法操作" width="380px" :close-on-click-modal="false">
      <p>请在班级详情页面进入课表模块。</p>
      <template #footer>
        <el-button type="primary" @click="goBack">知道了</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { reactive, ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { useTimetableStore } from '@/stores/timetable'
import { DAY_LABELS } from '@/utils/time'

const route = useRoute()
const router = useRouter()
const timetableStore = useTimetableStore()

const classId = computed(() => Number(route.query.class_id) || 0)
const isEdit = ref(false)
const submitting = ref(false)
const showPrereq = ref(false)
const dayLabels = DAY_LABELS

const entries = reactive([])

function addRow() {
  entries.push({
    day_of_week: '',
    period_start: '',
    period_end: '',
    course_name: '',
    teacher: '',
    room: '',
  })
}

function removeRow(idx) {
  entries.splice(idx, 1)
}

async function handleSubmit() {
  // 验证
  if (entries.length === 0) {
    ElMessage.warning('请至少添加一个课程条目')
    return
  }

  for (const [idx, row] of entries.entries()) {
    if (!row.day_of_week || !row.period_start || !row.period_end || !row.course_name) {
      ElMessage.warning(`第 ${idx + 1} 行：请填写完整的课程信息（星期、节次、课程名称）`)
      return
    }
    if (row.period_start > row.period_end) {
      ElMessage.warning(`第 ${idx + 1} 行：开始节次不能大于结束节次`)
      return
    }
  }

  if (!classId.value) {
    showPrereq.value = true
    return
  }

  submitting.value = true
  try {
    const payload = entries.map(e => ({
      day_of_week: Number(e.day_of_week),
      period_start: Number(e.period_start),
      period_end: Number(e.period_end),
      course_name: e.course_name,
      teacher: e.teacher || null,
      room: e.room || null,
    }))

    if (isEdit.value) {
      await timetableStore.updateClassTimetable(classId.value, payload)
      ElMessage.success('课表已更新，所有成员已重新继承')
    } else {
      await timetableStore.createClassTimetable(classId.value, payload)
      ElMessage.success('课表录入成功')
    }
    goBack()
  } catch (err) {
    ElMessage.error(err.message)
  } finally {
    submitting.value = false
  }
}

function goBack() {
  router.push({
    name: 'Timetable',
    query: { class_id: classId.value },
  })
}

// 页面加载时，获取已有课表（编辑模式）
onMounted(async () => {
  if (classId.value) {
    try {
      const data = await timetableStore.fetchClassTimetable(classId.value)
      if (data?.entries?.length > 0) {
        isEdit.value = true
        for (const e of data.entries) {
          entries.push({
            day_of_week: e.day_of_week,
            period_start: e.period_start,
            period_end: e.period_end,
            course_name: e.course_name,
            teacher: e.teacher || '',
            room: e.room || '',
          })
        }
      }
    } catch {
      // 无课表，新建模式
    }
  }
  // 如果没有 class_id，显示空列表
  if (!classId.value) {
    addRow()
  }
})
</script>
