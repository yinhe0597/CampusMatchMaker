import { defineStore } from 'pinia'
import { ref } from 'vue'
import { classApi } from '@/api/class'

export const useClassStore = defineStore('class', () => {
  // 状态
  const myClasses = ref([])
  const currentClass = ref(null)
  const members = ref([])
  const loading = ref(false)

  // 获取我的班级列表
  async function fetchMyClasses() {
    loading.value = true
    try {
      const res = await classApi.listMyClasses()
      if (res.code === 0) {
        myClasses.value = res.data || []
        return myClasses.value
      }
      throw new Error(res.message || '获取班级列表失败')
    } finally {
      loading.value = false
    }
  }

  // 获取班级详情
  async function fetchClassDetail(id) {
    const res = await classApi.getClassDetail(id)
    if (res.code === 0) {
      currentClass.value = res.data
      return res.data
    }
    throw new Error(res.message || '获取班级详情失败')
  }

  // 创建班级
  async function createClass(data) {
    const res = await classApi.createClass(data)
    if (res.code === 0) {
      await fetchMyClasses()
      return res.data
    }
    throw new Error(res.message || '创建班级失败')
  }

  // 通过邀请码查找班级
  async function lookupByCode(code) {
    const res = await classApi.lookupByCode(code)
    if (res.code === 0) {
      return res.data
    }
    throw new Error(res.message || '查找班级失败')
  }

  // 加入班级（先查后加）
  async function joinClass(inviteCode) {
    // 1. 查找班级
    const info = await lookupByCode(inviteCode)
    // 2. 加入班级
    const res = await classApi.joinClass(info.id, inviteCode)
    if (res.code === 0) {
      await fetchMyClasses()
      return res.data
    }
    throw new Error(res.message || '加入班级失败')
  }

  // 获取成员列表
  async function fetchMembers(classId, params = {}) {
    const res = await classApi.listMembers(classId, params)
    if (res.code === 0) {
      members.value = res.data?.list || []
      return {
        list: members.value,
        total: res.data?.total || 0,
      }
    }
    throw new Error(res.message || '获取成员列表失败')
  }

  // 移除成员
  async function removeMember(classId, userId) {
    const res = await classApi.removeMember(classId, userId)
    if (res.code === 0) {
      await fetchMembers(classId)
      return true
    }
    throw new Error(res.message || '移除成员失败')
  }

  // 清除当前班级
  function clearCurrent() {
    currentClass.value = null
    members.value = []
  }

  return {
    myClasses,
    currentClass,
    members,
    loading,
    fetchMyClasses,
    fetchClassDetail,
    createClass,
    lookupByCode,
    joinClass,
    fetchMembers,
    removeMember,
    clearCurrent,
  }
})
