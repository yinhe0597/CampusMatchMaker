<template>
  <div>
    <el-button size="small" @click="$router.push('/polls')">← 返回列表</el-button>

    <el-card style="margin-top: 12px;">
      <template #header>
        <div style="display: flex; justify-content: space-between; align-items: center;">
          <div>
            <span style="font-size: 18px; font-weight: bold;">{{ poll?.title }}</span>
            <el-tag :type="statusType" style="margin-left: 12px;" size="small">
              {{ statusLabel }}
            </el-tag>
          </div>
          <div v-if="isCreator">
            <el-button v-if="poll?.status === 'draft'" type="primary" size="small" @click="handleOpen">开启投票</el-button>
            <el-button v-if="poll?.status === 'open'" type="warning" size="small" @click="handleClose">关闭投票</el-button>
          </div>
        </div>
      </template>

      <div v-if="poll?.description" style="margin-bottom: 16px; color: #666;">
        {{ poll.description }}
      </div>

      <div style="font-size: 13px; color: #999; margin-bottom: 16px;">
        <span>截止时间：{{ poll?.deadline ? formatTime(poll.deadline) : '无' }}</span>
        <span style="margin-left: 16px;">投票人数：{{ poll?.voter_count || 0 }}</span>
        <span v-if="poll?.final_option_id" style="margin-left: 16px; color: #67c23a;">✓ 已确认最终时段</span>
      </div>
    </el-card>

    <!-- 选项列表 -->
    <el-card style="margin-top: 12px;">
      <template #header><span>时段选项</span></template>

      <div v-if="loading" style="text-align: center; padding: 20px;">
        <el-icon class="is-loading" :size="24"><Loading /></el-icon>
      </div>

      <template v-else-if="poll?.options?.length === 0">
        <el-empty description="暂无时段选项" />
      </template>

      <div v-else class="option-list">
        <div
          v-for="opt in poll?.options || []"
          :key="opt.option.id"
          class="option-card"
          :class="{ 'option-finalized': poll?.final_option_id === opt.option.id }"
        >
          <div class="option-info">
            <div class="option-time">
              <strong>{{ opt.option.slot_date }}</strong>
              <span style="margin-left: 8px;">{{ opt.option.slot_start_time?.substring(0, 5) }} - {{ opt.option.slot_end_time?.substring(0, 5) }}</span>
            </div>
            <div v-if="opt.option.is_recommended" style="font-size: 12px; color: #409eff; margin-top: 2px;">
              推荐参与率: {{ (opt.option.recommendation_rate * 100).toFixed(1) }}%
            </div>
            <div v-if="poll?.final_option_id === opt.option.id" class="finalized-badge">已确认</div>
          </div>

          <!-- 投票状态：进行中 -->
          <div v-if="poll?.status === 'open'" class="option-vote">
            <el-radio-group
              :model-value="getMyChoice(opt.option.id)"
              size="small"
              @change="(val) => handleVote(opt.option.id, val)"
            >
              <el-radio-button value="yes">可以</el-radio-button>
              <el-radio-button value="no">不行</el-radio-button>
              <el-radio-button value="maybe">待定</el-radio-button>
            </el-radio-group>
          </div>

          <!-- 投票结果：已关闭/已确认 -->
          <div v-if="poll?.status === 'closed' || poll?.status === 'finalized'" class="option-result">
            <el-progress
              :percentage="calcPercentage(opt)"
              :status="poll?.final_option_id === opt.option.id ? 'success' : undefined"
              :stroke-width="16"
            >
              <span style="font-size: 12px;">
                可: {{ opt.yes_count }} / 否: {{ opt.no_count }} / 待定: {{ opt.maybe_count }}
              </span>
            </el-progress>
          </div>
        </div>
      </div>
    </el-card>

    <!-- 确认最终时段（仅创建者可操作） -->
    <el-card v-if="isCreator && (poll?.status === 'closed' || poll?.status === 'open')" style="margin-top: 12px;">
      <template #header><span>确认最终时段</span></template>
      <div v-if="poll?.options?.length === 0">
        <el-empty description="暂无选项可确认" />
      </div>
      <div v-else style="display: flex; gap: 8px; flex-wrap: wrap;">
        <el-button
          v-for="opt in poll?.options || []"
          :key="'final-' + opt.option.id"
          :type="poll?.final_option_id === opt.option.id ? 'success' : 'default'"
          @click="handleFinalize(opt.option.id)"
        >
          {{ opt.option.slot_date?.substring(5) }} {{ opt.option.slot_start_time?.substring(0, 5) }}-
          {{ opt.option.slot_end_time?.substring(0, 5) }}
          (可:{{ opt.yes_count }})
        </el-button>
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Loading } from '@element-plus/icons-vue'
import { usePollStore } from '@/stores/poll'
import { useUserStore } from '@/stores/user'

const route = useRoute()
const router = useRouter()
const pollStore = usePollStore()
const userStore = useUserStore()

const poll = ref(null)
const loading = ref(false)

const isCreator = computed(() => poll.value?.creator_user_id === userStore.userInfo?.id)

const statusType = computed(() => {
  const map = { draft: 'info', open: 'success', closed: 'warning', finalized: 'primary' }
  return map[poll.value?.status] || 'info'
})

const statusLabel = computed(() => {
  const map = { draft: '草稿', open: '进行中', closed: '已关闭', finalized: '已确认' }
  return map[poll.value?.status] || poll.value?.status || ''
})

function getMyChoice(optionId) {
  const found = pollStore.myVotes?.find((v) => v.option_id === optionId)
  return found?.choice || undefined
}

function calcPercentage(opt) {
  const total = opt.yes_count + opt.no_count + opt.maybe_count
  if (total === 0) return 0
  return Math.round((opt.yes_count / total) * 100)
}

function formatTime(t) {
  if (!t) return ''
  return new Date(t).toLocaleString()
}

async function loadDetail() {
  loading.value = true
  try {
    const data = await pollStore.fetchPollDetail(route.params.id)
    poll.value = data
  } catch (e) {
    ElMessage.error(e.message || '获取详情失败')
  } finally {
    loading.value = false
  }
}

async function handleOpen() {
  try {
    await pollStore.openPoll(route.params.id)
    ElMessage.success('投票已开启')
    await loadDetail()
  } catch (e) {
    ElMessage.error(e.message || '开启失败')
  }
}

async function handleClose() {
  try {
    await ElMessageBox.confirm('确定关闭该投票？', '提示')
    await pollStore.closePoll(route.params.id)
    ElMessage.success('投票已关闭')
    await loadDetail()
  } catch (e) {
    if (e !== 'cancel') ElMessage.error(e.message || '关闭失败')
  }
}

async function handleVote(optionId, choice) {
  try {
    const votes = (pollStore.myVotes || [])
      .filter((v) => v.option_id !== optionId)
      .concat([{ option_id: optionId, choice }])
    // 确保每个选项都有投票 - 补全未投票的选项
    for (const opt of poll.value?.options || []) {
      if (!votes.find((v) => v.option_id === opt.option.id)) {
        votes.push({ option_id: opt.option.id, choice: 'maybe' })
      }
    }
    await pollStore.submitVote(route.params.id, votes)
    ElMessage.success('投票成功')
    await loadDetail()
  } catch (e) {
    ElMessage.error(e.message || '投票失败')
  }
}

async function handleFinalize(optionId) {
  try {
    await ElMessageBox.confirm('确定确认该时段为最终结果？', '提示')
    await pollStore.finalizePoll(route.params.id, optionId)
    ElMessage.success('已确认最终时段')
    await loadDetail()
  } catch (e) {
    if (e !== 'cancel') ElMessage.error(e.message || '确认失败')
  }
}

onMounted(loadDetail)
</script>

<style scoped>
.option-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.option-card {
  border: 1px solid #e4e7ed;
  border-radius: 8px;
  padding: 16px;
  transition: all 0.2s;
}

.option-card:hover {
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.08);
}

.option-finalized {
  border-color: #67c23a;
  background: #f0f9eb;
}

.option-info {
  position: relative;
}

.option-time {
  font-size: 15px;
}

.finalized-badge {
  position: absolute;
  top: 0;
  right: 0;
  background: #67c23a;
  color: #fff;
  font-size: 12px;
  padding: 2px 8px;
  border-radius: 4px;
}

.option-vote {
  margin-top: 12px;
}

.option-result {
  margin-top: 12px;
}
</style>
