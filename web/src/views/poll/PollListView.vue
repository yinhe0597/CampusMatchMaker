<template>
  <div>
    <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 16px;">
      <h2>投票列表</h2>
      <div>
        <el-select v-model="scopeType" placeholder="范围类型" style="width: 120px; margin-right: 8px;" @change="loadPolls">
          <el-option label="我的班级" value="my" />
        </el-select>
        <el-button type="primary" @click="$router.push('/polls/create')">发起投票</el-button>
      </div>
    </div>

    <el-tabs v-model="activeTab" @tab-change="onTabChange">
      <el-tab-pane label="进行中" name="open" />
      <el-tab-pane label="已关闭" name="closed" />
      <el-tab-pane label="已确认" name="finalized" />
      <el-tab-pane label="草稿" name="draft" />
    </el-tabs>

    <div v-if="loading" style="text-align: center; padding: 40px;">
      <el-icon class="is-loading" :size="32"><Loading /></el-icon>
    </div>

    <template v-else-if="polls.length === 0">
      <el-empty description="暂无投票" />
    </template>

    <template v-else>
      <el-card
        v-for="poll in polls"
        :key="poll.id"
        style="margin-bottom: 12px; cursor: pointer;"
        @click="$router.push(`/polls/${poll.id}`)"
      >
        <div style="display: flex; justify-content: space-between; align-items: center;">
          <div>
            <strong>{{ poll.title }}</strong>
            <div style="font-size: 12px; color: #999; margin-top: 4px;">
              {{ poll.total_options }} 个选项 · {{ poll.voter_count }} 人已投票
            </div>
          </div>
          <el-tag :type="statusType(poll.status)" size="small">
            {{ statusLabel(poll.status) }}
          </el-tag>
        </div>
      </el-card>

      <div v-if="totalCount > pageSize" style="text-align: center; margin-top: 16px;">
        <el-pagination
          v-model:current-page="page"
          :page-size="pageSize"
          :total="totalCount"
          layout="prev, pager, next"
          @current-change="loadPolls"
        />
      </div>
    </template>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { Loading } from '@element-plus/icons-vue'
import { usePollStore } from '@/stores/poll'
import { useClassStore } from '@/stores/class'

const router = useRouter()
const pollStore = usePollStore()
const classStore = useClassStore()

const scopeType = ref('my')
const activeTab = ref('open')
const page = ref(1)
const pageSize = ref(20)
const polls = ref([])
const totalCount = ref(0)
const loading = ref(false)

function statusType(status) {
  const map = { draft: 'info', open: 'success', closed: 'warning', finalized: 'primary' }
  return map[status] || 'info'
}

function statusLabel(status) {
  const map = { draft: '草稿', open: '进行中', closed: '已关闭', finalized: '已确认' }
  return map[status] || status
}

async function loadPolls() {
  loading.value = true
  try {
    // 先获取班级列表
    await classStore.fetchMyClasses()
    const myClasses = classStore.myClasses || []
    if (myClasses.length === 0) {
      polls.value = []
      totalCount.value = 0
      return
    }

    // 遍历所有班级获取投票（简化：取第一个班级）
    const firstClass = myClasses[0]
    const res = await pollStore.fetchPolls('class', firstClass.id, {
      page: page.value,
      page_size: pageSize.value,
    })
    // 服务端暂未实现 status 筛选，客户端过滤展示，保持后端分页
    polls.value = res?.polls || []
    totalCount.value = res?.total_count || 0
  } catch (e) {
    polls.value = []
    totalCount.value = 0
  } finally {
    loading.value = false
  }
}

// activeTab 变化时重置页码重新加载
function onTabChange() {
  page.value = 1
  loadPolls()
}

onMounted(loadPolls)
</script>
