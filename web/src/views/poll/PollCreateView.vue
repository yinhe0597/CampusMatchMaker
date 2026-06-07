<template>
  <div>
    <h2>发起时间投票</h2>
    <el-card>
      <el-form label-width="120px" style="max-width: 600px;">
        <el-form-item label="所属班级">
          <el-select v-model="form.scope_id" placeholder="选择班级" style="width: 100%;">
            <el-option
              v-for="cls in classStore.myClasses"
              :key="cls.id"
              :label="cls.name"
              :value="cls.id"
            />
          </el-select>
        </el-form-item>

        <el-form-item label="投票标题">
          <el-input v-model="form.title" placeholder="如：本周班会时间" />
        </el-form-item>

        <el-form-item label="投票说明">
          <el-input v-model="form.description" type="textarea" :rows="3" />
        </el-form-item>

        <el-form-item label="截止时间">
          <el-date-picker v-model="form.deadline" type="datetime" placeholder="选择截止时间" style="width: 100%;" />
        </el-form-item>

        <el-form-item label="自动推荐">
          <el-switch v-model="form.auto_recommend" active-text="系统根据课表自动推荐空闲时段" />
        </el-form-item>

        <template v-if="form.auto_recommend">
          <el-form-item label="推荐数量">
            <el-input-number v-model="timePref.max_recommendations" :min="1" :max="10" />
          </el-form-item>
          <el-form-item label="最短时长">
            <el-select v-model="timePref.min_duration_minutes">
              <el-option label="30 分钟" :value="30" />
              <el-option label="60 分钟" :value="60" />
              <el-option label="90 分钟" :value="90" />
              <el-option label="120 分钟" :value="120" />
            </el-select>
          </el-form-item>
          <el-form-item label="开始时间">
            <el-time-select v-model="dayStart" :max-time="dayEnd" start="08:00" step="01:00" end="22:00" placeholder="开始时间" />
          </el-form-item>
          <el-form-item label="结束时间">
            <el-time-select v-model="dayEnd" :min-time="dayStart" start="08:00" step="01:00" end="22:00" placeholder="结束时间" />
          </el-form-item>
        </template>

        <el-form-item>
          <el-button type="primary" :loading="submitting" @click="handleSubmit">
            创建投票{{ form.auto_recommend ? '（系统自动推荐空闲时段）' : '' }}
          </el-button>
          <el-button @click="$router.push('/polls')">取消</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup>
import { reactive, ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { usePollStore } from '@/stores/poll'
import { useClassStore } from '@/stores/class'

const router = useRouter()
const pollStore = usePollStore()
const classStore = useClassStore()

const submitting = ref(false)
const dayStart = ref('08:00')
const dayEnd = ref('22:00')

const form = reactive({
  title: '',
  description: '',
  scope_id: null,
  deadline: null,
  auto_recommend: false,
})

const timePref = reactive({
  day_start_hour: 8,
  day_end_hour: 22,
  min_duration_minutes: 60,
  max_recommendations: 3,
})

onMounted(async () => {
  await classStore.fetchMyClasses()
  if (classStore.myClasses.length > 0) {
    form.scope_id = classStore.myClasses[0].id
  }
})

function parseHour(timeStr) {
  if (!timeStr) return 8
  const parts = timeStr.split(':')
  return parseInt(parts[0]) || 8
}

async function handleSubmit() {
  if (!form.title) {
    ElMessage.warning('请输入投票标题')
    return
  }
  if (!form.scope_id) {
    ElMessage.warning('请选择班级')
    return
  }

  submitting.value = true
  try {
    const data = {
      title: form.title,
      description: form.description || undefined,
      scope_type: 'class',
      scope_id: form.scope_id,
      deadline: form.deadline || undefined,
      auto_recommend: form.auto_recommend,
    }

    if (form.auto_recommend) {
      data.time_preference = {
        day_start_hour: parseHour(dayStart.value),
        day_end_hour: parseHour(dayEnd.value),
        min_duration_minutes: timePref.min_duration_minutes,
        max_recommendations: timePref.max_recommendations,
      }
    }

    const result = await pollStore.createPoll(data)
    ElMessage.success(`投票创建成功${result.options_created > 0 ? `，已推荐 ${result.options_created} 个时段` : ''}`)
    router.push(`/polls/${result.poll_id}`)
  } catch (e) {
    ElMessage.error(e.message || '创建失败')
  } finally {
    submitting.value = false
  }
}
</script>
