<template>
  <div>
    <!-- 返回按钮 -->
    <el-button text :icon="'ArrowLeft'" @click="goBack" style="margin-bottom: 12px;">
      返回班级列表
    </el-button>

    <!-- 加载状态 -->
    <div v-if="loading" style="text-align: center; padding: 80px 0;">
      <el-skeleton :rows="4" animated />
    </div>

    <template v-else-if="detail">
      <!-- 基本信息卡片 -->
      <el-card shadow="never" style="margin-bottom: 16px;">
        <div style="display: flex; justify-content: space-between; align-items: flex-start;">
          <div>
            <h2 style="margin: 0 0 4px 0;">{{ detail.name }}</h2>
            <p style="margin: 0; color: #909399; font-size: 14px;">
              {{ detail.grade }}
              <template v-if="detail.department"> · {{ detail.department }}</template>
              · {{ detail.member_count }} 人
            </p>
          </div>
          <el-tag :type="myRoleTag" size="small">{{ myRoleLabel }}</el-tag>
        </div>

        <!-- 邀请码（仅管理员可见） -->
        <div v-if="detail.invite_code" style="margin-top: 16px; padding: 12px; background: #f5f7fa; border-radius: 6px; display: flex; align-items: center; gap: 12px;">
          <span style="font-size: 14px; color: #606266;">邀请码：</span>
          <span style="font-family: monospace; font-size: 20px; font-weight: bold; letter-spacing: 4px; color: #409eff;">
            {{ detail.invite_code }}
          </span>
          <el-button text type="primary" size="small" @click="copyInviteCode">复制</el-button>
        </div>
      </el-card>

      <!-- Tabs -->
      <el-card shadow="never">
        <el-tabs v-model="activeTab">
          <el-tab-pane label="成员列表" name="members">
            <!-- 成员表格 -->
            <el-table :data="classStore.members" stripe style="width: 100%;" v-loading="memberLoading">
              <el-table-column label="用户" min-width="160">
                <template #default="{ row }">
                  <div style="display: flex; align-items: center; gap: 8px;">
                    <el-avatar :size="32" icon="User" />
                    <span>{{ row.nickname || '用户' + row.user_id }}</span>
                  </div>
                </template>
              </el-table-column>
              <el-table-column prop="role" label="角色" width="100">
                <template #default="{ row }">
                  <el-tag :type="roleTagType(row.role)" size="small">
                    {{ roleLabel(row.role) }}
                  </el-tag>
                </template>
              </el-table-column>
              <el-table-column prop="joined_at" label="加入时间" width="180" />
              <el-table-column label="操作" width="120" v-if="canManage">
                <template #default="{ row }">
                  <el-button
                    v-if="row.user_id !== userStore.userInfo?.id && row.role !== 'owner'"
                    type="danger"
                    text
                    size="small"
                    @click="handleRemove(row)"
                  >
                    移除
                  </el-button>
                </template>
              </el-table-column>
            </el-table>

            <!-- 分页 -->
            <div style="display: flex; justify-content: flex-end; margin-top: 16px;">
              <el-pagination
                v-if="memberTotal > 20"
                v-model:current-page="memberPage"
                :page-size="20"
                :total="memberTotal"
                layout="prev, pager, next"
                small
                @current-change="loadMembers"
              />
            </div>
          </el-tab-pane>
          <el-tab-pane label="班级课表" name="timetable">
            <el-empty description="暂无课表">
              <el-button v-if="canManage" type="primary" @click="goToTimetable">
                {{ detail?.timetable_status === 1 ? '查看/编辑课表' : '录入课表' }}
              </el-button>
            </el-empty>
          </el-tab-pane>
          <el-tab-pane label="投票列表" name="polls">
            <div v-if="pollLoading" style="text-align: center; padding: 20px;">
              <el-icon class="is-loading" :size="24"><Loading /></el-icon>
            </div>
            <el-empty v-else-if="classPolls.length === 0" description="暂无投票" />
            <template v-else>
              <el-card
                v-for="poll in classPolls"
                :key="poll.id"
                style="margin-bottom: 8px; cursor: pointer;"
                @click="$router.push(`/polls/${poll.id}`)"
              >
                <div style="display: flex; justify-content: space-between; align-items: center;">
                  <div>
                    <strong>{{ poll.title }}</strong>
                    <div style="font-size: 12px; color: #999; margin-top: 2px;">
                      {{ poll.total_options }} 个选项 · {{ poll.voter_count }} 人已投票
                    </div>
                  </div>
                  <el-tag :type="pollStatusType(poll.status)" size="small">
                    {{ pollStatusLabel(poll.status) }}
                  </el-tag>
                </div>
              </el-card>
            </template>
          </el-tab-pane>
        </el-tabs>
      </el-card>
    </template>

    <!-- 错误状态 -->
    <el-result v-else status="error" title="班级不存在" sub-title="请检查班级ID是否正确">
      <template #extra>
        <el-button type="primary" @click="goBack">返回班级列表</el-button>
      </template>
    </el-result>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Loading } from '@element-plus/icons-vue'
import { useClassStore } from '@/stores/class'
import { useUserStore } from '@/stores/user'
import { usePollStore } from '@/stores/poll'

const route = useRoute()
const router = useRouter()
const classStore = useClassStore()
const userStore = useUserStore()
const pollStore = usePollStore()

const loading = ref(true)
const memberLoading = ref(false)
const activeTab = ref('members')
const memberPage = ref(1)
const memberTotal = ref(0)

const detail = computed(() => classStore.currentClass)

// 投票相关
const classPolls = ref([])
const pollLoading = ref(false)

function pollStatusType(s) {
  const map = { draft: 'info', open: 'success', closed: 'warning', finalized: 'primary' }
  return map[s] || 'info'
}
function pollStatusLabel(s) {
  const map = { draft: '草稿', open: '进行中', closed: '已关闭', finalized: '已确认' }
  return map[s] || s
}

async function loadClassPolls() {
  const classId = Number(route.params.id)
  if (!classId) return
  pollLoading.value = true
  try {
    const res = await pollStore.fetchPolls('class', classId, { page: 1, page_size: 50 })
    classPolls.value = res?.polls || []
  } catch {
    classPolls.value = []
  } finally {
    pollLoading.value = false
  }
}

// 切换到投票tab时加载
watch(activeTab, (val) => {
  if (val === 'polls') loadClassPolls()
})
const canManage = computed(() => {
  const role = detail.value?.my_role
  return role === 'owner' || role === 'admin'
})
const myRoleLabel = computed(() => roleLabel(detail.value?.my_role))
const myRoleTag = computed(() => roleTagType(detail.value?.my_role))

// 加载班级详情
async function loadDetail() {
  const classId = Number(route.params.id)
  if (!classId) {
    loading.value = false
    return
  }

  loading.value = true
  try {
    await classStore.fetchClassDetail(classId)
    await loadMembers()
  } catch {
    classStore.clearCurrent()
    ElMessage.error('获取班级详情失败')
  } finally {
    loading.value = false
  }
}

// 加载成员列表
async function loadMembers() {
  const classId = Number(route.params.id)
  memberLoading.value = true
  try {
    const result = await classStore.fetchMembers(classId, {
      page: memberPage.value,
      page_size: 20,
    })
    memberTotal.value = result.total
  } catch {
    ElMessage.error('获取成员列表失败')
  } finally {
    memberLoading.value = false
  }
}

// 移除成员
async function handleRemove(row) {
  try {
    await ElMessageBox.confirm(
      `确定要移除成员「${row.nickname || '用户' + row.user_id}」吗？`,
      '确认移除',
      { confirmButtonText: '移除', cancelButtonText: '取消', type: 'warning' }
    )
    await classStore.removeMember(Number(route.params.id), row.user_id)
    ElMessage.success('已移除')
  } catch {
    // 取消操作不处理
  }
}

// 复制邀请码
async function copyInviteCode() {
  try {
    await navigator.clipboard.writeText(detail.value?.invite_code || '')
    ElMessage.success('邀请码已复制')
  } catch {
    ElMessage.warning('复制失败，请手动复制')
  }
}

// 工具函数
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
function goBack() {
  router.push({ name: 'ClassList' })
}

function goToTimetable() {
  router.push({
    name: 'Timetable',
    query: {
      class_id: route.params.id,
      role: detail.value?.my_role
    }
  })
}

onMounted(loadDetail)
</script>
