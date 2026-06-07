import { defineStore } from 'pinia'
import { ref } from 'vue'
import { timetableApi } from '@/api/timetable'

export const useTimetableStore = defineStore('timetable', () => {
  // 班级公共课表
  const classTimetable = ref(null)
  // 个人课表
  const personalTimetable = ref(null)
  // 纠错列表
  const corrections = ref([])
  // 加载状态
  const loading = ref(false)

  // ===== 班级公共课表 =====

  async function fetchClassTimetable(classId) {
    loading.value = true
    try {
      const res = await timetableApi.getClassTimetable(classId)
      if (res.code === 0) {
        classTimetable.value = res.data
        return res.data
      }
      throw new Error(res.message || '获取课表失败')
    } finally {
      loading.value = false
    }
  }

  async function createClassTimetable(classId, entries) {
    const res = await timetableApi.createClassTimetable(classId, { entries })
    if (res.code === 0) {
      await fetchClassTimetable(classId)
      return res.data
    }
    throw new Error(res.message || '录入课表失败')
  }

  async function updateClassTimetable(classId, entries) {
    const res = await timetableApi.updateClassTimetable(classId, { entries })
    if (res.code === 0) {
      await fetchClassTimetable(classId)
      return res.data
    }
    throw new Error(res.message || '更新课表失败')
  }

  // ===== 个人课表 =====

  async function fetchPersonalTimetable(classId) {
    loading.value = true
    try {
      const res = await timetableApi.getPersonalTimetable(classId)
      if (res.code === 0) {
        personalTimetable.value = res.data
        return res.data
      }
      throw new Error(res.message || '获取个人课表失败')
    } finally {
      loading.value = false
    }
  }

  async function createPersonalTimetable(data) {
    const res = await timetableApi.createPersonalTimetable(data)
    if (res.code === 0) {
      await fetchPersonalTimetable(data.class_id)
      return res.data
    }
    throw new Error(res.message || '添加课程失败')
  }

  async function updatePersonalTimetable(id, data) {
    const res = await timetableApi.updatePersonalTimetable(id, data)
    if (res.code === 0) {
      return true
    }
    throw new Error(res.message || '更新失败')
  }

  async function deletePersonalTimetable(id, classId) {
    const res = await timetableApi.deletePersonalTimetable(id)
    if (res.code === 0) {
      await fetchPersonalTimetable(classId)
      return true
    }
    throw new Error(res.message || '删除失败')
  }

  // ===== 纠错 =====

  async function fetchCorrections(params) {
    const res = await timetableApi.listCorrections(params)
    if (res.code === 0) {
      corrections.value = res.data?.list || []
      return res.data
    }
    throw new Error(res.message || '获取纠错列表失败')
  }

  async function createCorrection(data) {
    const res = await timetableApi.createCorrection(data)
    if (res.code === 0) {
      return res.data
    }
    throw new Error(res.message || '提交纠错失败')
  }

  async function reviewCorrection(id, action) {
    const res = await timetableApi.reviewCorrection(id, { action })
    if (res.code === 0) {
      return true
    }
    throw new Error(res.message || '处理纠错失败')
  }

  function clear() {
    classTimetable.value = null
    personalTimetable.value = null
    corrections.value = []
  }

  return {
    classTimetable,
    personalTimetable,
    corrections,
    loading,
    fetchClassTimetable,
    createClassTimetable,
    updateClassTimetable,
    fetchPersonalTimetable,
    createPersonalTimetable,
    updatePersonalTimetable,
    deletePersonalTimetable,
    fetchCorrections,
    createCorrection,
    reviewCorrection,
    clear,
  }
})
