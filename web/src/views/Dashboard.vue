<template>
  <div class="min-h-screen bg-gray-50 p-4 md:p-8">
    <div class="max-w-5xl mx-auto">
      <div class="flex items-center justify-between mb-6 gap-4">
        <h1 class="text-xl md:text-3xl font-bold text-gray-800 truncate">SMS Dashboard</h1>
        <div class="flex items-center gap-2 flex-shrink-0">
            <button @click="fetchSMS" class="px-3 py-1.5 md:px-4 md:py-2 bg-blue-500 text-white rounded hover:bg-blue-600 transition text-sm whitespace-nowrap">Refresh</button>
            <button @click="showChangePassword = true" class="px-3 py-1.5 md:px-4 md:py-2 bg-gray-500 text-white rounded hover:bg-gray-600 transition text-sm whitespace-nowrap">Password</button>
            <button @click="logout" class="px-3 py-1.5 md:px-4 md:py-2 bg-red-500 text-white rounded hover:bg-red-600 transition text-sm whitespace-nowrap">Logout</button>
        </div>
      </div>

      <!-- Change Password Modal -->
      <div v-if="showChangePassword" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
        <div class="bg-white rounded-lg shadow-xl w-full max-w-md p-6">
          <h2 class="text-xl font-bold mb-4">Change Password</h2>
          <div class="space-y-4">
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Old Password</label>
              <input v-model="passwordForm.oldPassword" type="password" class="w-full px-3 py-2 border rounded focus:ring-2 focus:ring-blue-500 outline-none" placeholder="Enter old password">
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">New Password</label>
              <input v-model="passwordForm.newPassword" type="password" class="w-full px-3 py-2 border rounded focus:ring-2 focus:ring-blue-500 outline-none" placeholder="Min 6 characters">
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Confirm New Password</label>
              <input v-model="passwordForm.confirmPassword" type="password" class="w-full px-3 py-2 border rounded focus:ring-2 focus:ring-blue-500 outline-none" placeholder="Repeat new password">
            </div>
          </div>
          <div class="mt-6 flex justify-end gap-3">
            <button @click="closePasswordModal" class="px-4 py-2 text-gray-600 hover:bg-gray-100 rounded transition">Cancel</button>
            <button @click="changePassword" :disabled="isChanging" class="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700 transition disabled:opacity-50">
              {{ isChanging ? 'Saving...' : 'Save Changes' }}
            </button>
          </div>
        </div>
      </div>
      
      <!-- Desktop Table View -->
      <div class="hidden md:block bg-white shadow rounded-lg overflow-hidden border border-gray-200">
        <table class="min-w-full divide-y divide-gray-200 table-fixed">
          <thead class="bg-gray-50">
            <tr>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider w-16">ID</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider w-48">Time</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Content</th>
              <th class="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider w-24">Actions</th>
            </tr>
          </thead>
          <tbody class="bg-white divide-y divide-gray-200">
            <tr v-for="sms in smsList" :key="sms.id" class="hover:bg-gray-50 transition">
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{{ sms.id }}</td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{{ sms.sendTime }}</td>
              <td class="px-6 py-4 text-sm text-gray-900 break-words leading-relaxed">
                {{ sms.content }}
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                <button @click="deleteSMS(sms.id)" class="text-red-600 hover:text-red-900 transition">Delete</button>
              </td>
            </tr>
            <tr v-if="smsList.length === 0">
                <td colspan="4" class="px-6 py-12 text-center text-gray-500">No SMS messages found</td>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- Mobile Card View -->
      <div class="md:hidden space-y-4">
        <div v-for="sms in smsList" :key="sms.id" class="bg-white shadow rounded-lg border border-gray-200 p-4 relative">
          <div class="flex justify-between items-start mb-2 border-b pb-2">
            <span class="text-xs font-bold text-gray-400">ID: {{ sms.id }}</span>
            <span class="text-xs text-gray-500">{{ sms.sendTime }}</span>
          </div>
          <div class="text-sm text-gray-900 break-words mb-4 leading-relaxed whitespace-pre-wrap">
            {{ sms.content }}
          </div>
          <div class="flex justify-end">
            <button @click="deleteSMS(sms.id)" class="px-3 py-1 text-xs font-medium text-red-600 border border-red-600 rounded hover:bg-red-50 transition">Delete</button>
          </div>
        </div>
        <div v-if="smsList.length === 0" class="bg-white shadow rounded-lg border border-gray-200 p-12 text-center text-gray-500">
          No SMS messages found
        </div>
      </div>

      <!-- Pagination -->
      <div class="mt-6 flex flex-col md:flex-row justify-between items-center bg-white p-4 shadow rounded-lg border border-gray-200 gap-4">
        <div class="text-sm text-gray-500 order-2 md:order-1">
          Showing <span class="font-medium">{{ (pagination.page - 1) * pagination.pageSize + 1 }}</span> to 
          <span class="font-medium">{{ Math.min(pagination.page * pagination.pageSize, pagination.total) }}</span> of 
          <span class="font-medium">{{ pagination.total }}</span> results
        </div>
        <div class="flex space-x-2 order-1 md:order-2 w-full md:w-auto">
          <button 
            @click="changePage(pagination.page - 1)" 
            :disabled="pagination.page <= 1"
            class="flex-1 md:flex-none px-3 py-2 md:py-1 border rounded text-sm transition"
            :class="pagination.page <= 1 ? 'bg-gray-100 text-gray-400 cursor-not-allowed' : 'bg-white text-gray-700 hover:bg-gray-50'"
          >
            Previous
          </button>
          <button 
            @click="changePage(pagination.page + 1)" 
            :disabled="pagination.page * pagination.pageSize >= pagination.total"
            class="flex-1 md:flex-none px-3 py-2 md:py-1 border rounded text-sm transition"
            :class="pagination.page * pagination.pageSize >= pagination.total ? 'bg-gray-100 text-gray-400 cursor-not-allowed' : 'bg-white text-gray-700 hover:bg-gray-50'"
          >
            Next
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import axios from 'axios'
import { useRouter } from 'vue-router'

interface SMS {
  id: number
  sendTime: string
  content: string
}

interface Pagination {
  page: number
  pageSize: number
  total: number
}

const smsList = ref<SMS[]>([])
const pagination = ref<Pagination>({
  page: 1,
  pageSize: 10,
  total: 0
})
const router = useRouter()

// Change Password State
const showChangePassword = ref(false)
const isChanging = ref(false)
const passwordForm = ref({
  oldPassword: '',
  newPassword: '',
  confirmPassword: ''
})

const closePasswordModal = () => {
  showChangePassword.value = false
  passwordForm.value = {
    oldPassword: '',
    newPassword: '',
    confirmPassword: ''
  }
}

const changePassword = async () => {
  if (passwordForm.value.newPassword !== passwordForm.value.confirmPassword) {
    alert('New passwords do not match')
    return
  }
  if (passwordForm.value.newPassword.length < 6) {
    alert('New password must be at least 6 characters')
    return
  }

  isChanging.value = true
  try {
    const token = localStorage.getItem('token')
    await axios.post('/api/change-password', {
      oldPassword: passwordForm.value.oldPassword,
      newPassword: passwordForm.value.newPassword
    }, {
      headers: { Authorization: `Bearer ${token}` }
    })
    alert('Password changed successfully! Please login again.')
    logout()
  } catch (e) {
    console.error(e)
    const msg = axios.isAxiosError(e) ? e.response?.data?.error : 'Failed to change password'
    alert(msg || 'Failed to change password')
  } finally {
    isChanging.value = false
  }
}

const fetchSMS = async () => {
  try {
    const token = localStorage.getItem('token')
    if (!token) {
        router.push('/login')
        return
    }
    const res = await axios.get('/api/sms/list', {
      params: {
        page: pagination.value.page,
        pageSize: pagination.value.pageSize
      },
      headers: { Authorization: `Bearer ${token}` }
    })
    smsList.value = res.data.data
    pagination.value = res.data.pagination
  } catch (e) {
    console.error(e)
    if (axios.isAxiosError(e) && e.response?.status === 401) {
        logout()
    }
  }
}

const deleteSMS = async (id: number) => {
  if (!confirm('Are you sure you want to delete this SMS?')) return
  
  try {
    const token = localStorage.getItem('token')
    await axios.delete(`/api/sms/${id}`, {
      headers: { Authorization: `Bearer ${token}` }
    })
    fetchSMS()
  } catch (e) {
    console.error(e)
    alert('Failed to delete SMS')
  }
}

const changePage = (newPage: number) => {
  pagination.value.page = newPage
  fetchSMS()
}

const logout = () => {
  localStorage.removeItem('token')
  router.push('/login')
}

onMounted(() => {
  fetchSMS()
})
</script>
