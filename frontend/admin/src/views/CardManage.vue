<template>
  <el-card>
    <div class="toolbar" style="margin-bottom: 20px;">
      <el-button type="primary" @click="showDialog()">添加卡片</el-button>
      <el-input v-model="searchForm.keyword" placeholder="搜索卡号/设备号" style="width: 300px; margin-left: 10px;" @change="loadCards" />
      <el-select v-model="searchForm.status" placeholder="状态" style="width: 120px; margin-left: 10px;" @change="loadCards">
        <el-option label="全部" :value="0" />
        <el-option label="正常" :value="1" />
        <el-option label="即将到期" :value="2" />
        <el-option label="已过期" :value="3" />
      </el-select>
    </div>
    
    <el-table :data="cards" border v-loading="loading">
      <el-table-column prop="card_no" label="卡号" width="200" />
      <el-table-column prop="device_no" label="设备号" width="180" />
      <el-table-column prop="operator" label="运营商" width="100" />
      <el-table-column prop="start_date" label="开始日期" width="120" />
      <el-table-column prop="expire_date" label="到期日期" width="120" />
      <el-table-column prop="status_text" label="状态" width="100">
        <template #default="{ row }">
          <el-tag :type="row.status === 1 ? 'success' : row.status === 2 ? 'warning' : 'danger'">
            {{ row.status_text }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="remark" label="备注" show-overflow-tooltip />
      <el-table-column label="操作" width="180" fixed="right">
        <template #default="{ row }">
          <el-button link type="primary" @click="showDialog(row)">编辑</el-button>
          <el-button link type="danger" @click="handleDelete(row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>
    
    <el-pagination
      v-model:current-page="searchForm.page"
      v-model:page-size="searchForm.page_size"
      :total="total"
      @current-change="loadCards"
      style="margin-top: 20px; justify-content: center"
    />

    <!-- 添加/编辑对话框 -->
    <el-dialog v-model="dialogVisible" :title="form.id ? '编辑卡片' : '添加卡片'" width="600px">
      <el-form :model="form" label-width="100px">
        <el-form-item label="卡号">
          <el-input v-model="form.card_no" placeholder="请输入ICCID" />
        </el-form-item>
        <el-form-item label="设备号">
          <el-input v-model="form.device_no" placeholder="请输入IMEI" />
        </el-form-item>
        <el-form-item label="运营商">
          <el-select v-model="form.operator" style="width: 100%">
            <el-option label="中国移动" value="中国移动" />
            <el-option label="中国联通" value="中国联通" />
            <el-option label="中国电信" value="中国电信" />
          </el-select>
        </el-form-item>
        <el-form-item label="套餐类型">
          <el-input v-model="form.package_type" placeholder="如：年卡" />
        </el-form-item>
        <el-form-item label="开始日期">
          <el-date-picker v-model="form.start_date" type="date" style="width: 100%" value-format="YYYY-MM-DD" />
        </el-form-item>
        <el-form-item label="到期日期">
          <el-date-picker v-model="form.expire_date" type="date" style="width: 100%" value-format="YYYY-MM-DD" />
        </el-form-item>
        <el-form-item label="备注">
          <el-input v-model="form.remark" type="textarea" :rows="3" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSave" :loading="saving">保存</el-button>
      </template>
    </el-dialog>
  </el-card>
</template>

<script>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { getCardList, createCard, updateCard, deleteCard } from '@/api/admin'

export default {
  name: 'CardManage',
  setup() {
    const loading = ref(false)
    const saving = ref(false)
    const dialogVisible = ref(false)
    const cards = ref([])
    const total = ref(0)
    
    const searchForm = reactive({
      page: 1,
      page_size: 20,
      status: 0,
      keyword: ''
    })

    const form = reactive({
      id: null,
      card_no: '',
      device_no: '',
      operator: '中国移动',
      package_type: '年卡',
      start_date: '',
      expire_date: '',
      remark: ''
    })

    const loadCards = async () => {
      loading.value = true
      try {
        const data = await getCardList(searchForm)
        cards.value = data.list || []
        total.value = data.total || 0
      } catch (error) {
        console.error('加载卡片列表失败:', error)
      } finally {
        loading.value = false
      }
    }

    const showDialog = (row = null) => {
      if (row) {
        Object.assign(form, row)
      } else {
        Object.assign(form, {
          id: null,
          card_no: '',
          device_no: '',
          operator: '中国移动',
          package_type: '年卡',
          start_date: '',
          expire_date: '',
          remark: ''
        })
      }
      dialogVisible.value = true
    }

    const handleSave = async () => {
      saving.value = true
      try {
        if (form.id) {
          await updateCard(form.id, form)
          ElMessage.success('更新成功')
        } else {
          await createCard(form)
          ElMessage.success('添加成功')
        }
        dialogVisible.value = false
        loadCards()
      } catch (error) {
        console.error('保存失败:', error)
      } finally {
        saving.value = false
      }
    }

    const handleDelete = async (row) => {
      try {
        await ElMessageBox.confirm('确定要删除这张卡片吗？', '提示', { type: 'warning' })
        await deleteCard(row.id)
        ElMessage.success('删除成功')
        loadCards()
      } catch (error) {
        if (error !== 'cancel') {
          console.error('删除失败:', error)
        }
      }
    }

    onMounted(() => {
      loadCards()
    })

    return {
      loading,
      saving,
      dialogVisible,
      cards,
      total,
      searchForm,
      form,
      loadCards,
      showDialog,
      handleSave,
      handleDelete
    }
  }
}
</script>
