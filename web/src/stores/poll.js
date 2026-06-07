import { defineStore } from 'pinia'
import { ref } from 'vue'
import { pollApi } from '@/api/poll'

export const usePollStore = defineStore('poll', () => {
  const polls = ref([])
  const totalCount = ref(0)
  const currentPoll = ref(null)
  const myVotes = ref([])
  const loading = ref(false)

  // 获取投票列表
  async function fetchPolls(scopeType, scopeId, params = {}) {
    loading.value = true
    try {
      const res = await pollApi.listPolls({
        scope_type: scopeType,
        scope_id: scopeId,
        ...params,
      })
      if (res.code === 0) {
        polls.value = res.data?.polls || []
        totalCount.value = res.data?.total_count || 0
        return res.data
      }
      throw new Error(res.message || '获取投票列表失败')
    } finally {
      loading.value = false
    }
  }

  // 获取投票详情
  async function fetchPollDetail(id) {
    const res = await pollApi.getPollDetail(id)
    if (res.code === 0) {
      currentPoll.value = res.data
      myVotes.value = res.data?.my_votes || []
      return res.data
    }
    throw new Error(res.message || '获取投票详情失败')
  }

  // 创建投票
  async function createPoll(data) {
    const res = await pollApi.createPoll(data)
    if (res.code === 0) {
      return res.data
    }
    throw new Error(res.message || '创建投票失败')
  }

  // 开启投票
  async function openPoll(id) {
    const res = await pollApi.openPoll(id)
    if (res.code === 0) {
      if (currentPoll.value && currentPoll.value.id === id) {
        currentPoll.value.status = 'open'
      }
      return true
    }
    throw new Error(res.message || '开启投票失败')
  }

  // 关闭投票
  async function closePoll(id) {
    const res = await pollApi.closePoll(id)
    if (res.code === 0) {
      if (currentPoll.value && currentPoll.value.id === id) {
        currentPoll.value.status = 'closed'
      }
      return true
    }
    throw new Error(res.message || '关闭投票失败')
  }

  // 提交投票
  async function submitVote(id, votes) {
    const res = await pollApi.submitVote(id, { votes })
    if (res.code === 0) {
      // 更新本地 myVotes
      myVotes.value = votes.map((v) => ({ option_id: v.option_id, choice: v.choice }))
      return res.data
    }
    throw new Error(res.message || '提交投票失败')
  }

  // 获取结果
  async function fetchResults(id) {
    const res = await pollApi.getResults(id)
    if (res.code === 0) {
      return res.data?.items || []
    }
    throw new Error(res.message || '获取结果失败')
  }

  // 确认最终时段
  async function finalizePoll(id, optionId) {
    const res = await pollApi.finalizePoll(id, { final_option_id: optionId })
    if (res.code === 0) {
      if (currentPoll.value && currentPoll.value.id === id) {
        currentPoll.value.status = 'finalized'
        currentPoll.value.final_option_id = optionId
      }
      return true
    }
    throw new Error(res.message || '确认时段失败')
  }

  function clearCurrent() {
    currentPoll.value = null
    myVotes.value = []
  }

  return {
    polls,
    totalCount,
    currentPoll,
    myVotes,
    loading,
    fetchPolls,
    fetchPollDetail,
    createPoll,
    openPoll,
    closePoll,
    submitVote,
    fetchResults,
    finalizePoll,
    clearCurrent,
  }
})
