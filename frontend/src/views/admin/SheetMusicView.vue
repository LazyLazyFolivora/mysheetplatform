<template>
  <div class="sheet-manage">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>乐谱列表</span>
          <div class="header-actions">
            <el-input
              v-model="keyword"
              placeholder="搜索乐谱"
              style="width: 200px; margin-right: 16px"
              clearable
              @keyup.enter="handleSearch"
              @clear="handleSearch"
            />
            <el-button type="primary" @click="openCreate">添加乐谱</el-button>
          </div>
        </div>
      </template>

      <el-table :data="sheets" v-loading="loading" style="width: 100%">
        <el-table-column prop="id" label="ID" width="70" />
        <el-table-column prop="title" label="标题" min-width="140" />
        <el-table-column prop="composer" label="作曲家" width="120" />
        <el-table-column prop="description" label="描述" min-width="160" show-overflow-tooltip />
        <el-table-column prop="download_points" label="下载积分" width="90">
          <template #default="{ row }">{{ row.download_points || 0 }} 积分</template>
        </el-table-column>
        <el-table-column prop="price" label="价格" width="90">
          <template #default="{ row }">￥{{ (row.price || 0).toFixed(2) }}</template>
        </el-table-column>
        <el-table-column prop="view_count" label="浏览次数" width="90" />
        <el-table-column label="创建时间" width="170">
          <template #default="{ row }">{{ formatTime(row.created_at) }}</template>
        </el-table-column>
        <el-table-column label="操作" width="160" fixed="right">
          <template #default="{ row }">
            <el-button-group>
              <el-button type="primary" size="small" @click="openEdit(row)">编辑</el-button>
              <el-button type="danger" size="small" @click="handleDelete(row)">删除</el-button>
            </el-button-group>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination">
        <el-pagination
          :current-page="page"
          :page-size="pageSize"
          :total="total"
          layout="total, prev, pager, next"
          @current-change="handlePageChange"
        />
      </div>
    </el-card>

    <!-- 添加/编辑乐谱对话框 -->
    <el-dialog
      v-model="dialogVisible"
      :title="editingId ? '编辑乐谱' : '添加乐谱'"
      width="560px"
      :close-on-click-modal="false"
    >
      <el-form ref="formRef" :model="form" label-width="90px" v-loading="detailLoading">
        <el-form-item label="标题" required>
          <el-input v-model="form.title" placeholder="请输入乐谱标题" />
        </el-form-item>
        <el-form-item label="曲名">
          <el-input v-model="form.track_name" placeholder="例如：XX动漫主题曲" />
        </el-form-item>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="作曲">
              <el-input v-model="form.composer" placeholder="请输入作曲" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="编曲">
              <el-input v-model="form.arranger" placeholder="请输入编曲" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item label="制谱">
          <el-input v-model="form.transcriber" placeholder="请输入扒谱/制谱者" />
        </el-form-item>
        <el-form-item label="描述">
          <el-input v-model="form.description" type="textarea" :rows="3" placeholder="请输入乐谱描述" />
        </el-form-item>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="积分">
              <el-input-number v-model="form.download_points" :min="0" style="width: 100%" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="价格">
              <el-input-number v-model="form.price" :min="0" :precision="2" style="width: 100%" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item label="标签">
          <el-select
            v-model="form.tags"
            multiple
            filterable
            allow-create
            default-first-option
            placeholder="请选择或创建标签"
            style="width: 100%"
          >
            <el-option v-for="item in allTags" :key="item.id" :label="item.name" :value="item.name" />
          </el-select>
        </el-form-item>
        <el-form-item label="外部链接">
          <el-input v-model="form.external_link" placeholder="请输入外部链接" />
        </el-form-item>

        <el-form-item label="免费版乐谱">
          <div class="file-picker">
            <el-button size="small" @click="freePdfInput?.click()">
              {{ editingId ? '替换免费版' : '选择免费版 PDF' }}
            </el-button>
            <span v-if="freePdfFile" class="file-name">{{ freePdfFile.name }}</span>
            <span v-else-if="currentFreeFile" class="file-current">当前：{{ currentFreeFile.file_name }}</span>
            <input ref="freePdfInput" type="file" accept=".pdf" style="display: none" @change="onPickFreePdf" />
          </div>
          <div class="el-upload__tip">免费版 PDF 文件，且不超过 50MB</div>
        </el-form-item>

        <el-form-item label="付费版乐谱">
          <div class="file-picker">
            <el-button size="small" @click="paidPdfInput?.click()">
              {{ editingId ? '替换付费版' : '选择付费版 PDF' }}
            </el-button>
            <span v-if="paidPdfFile" class="file-name">{{ paidPdfFile.name }}</span>
            <span v-else-if="currentPaidFile" class="file-current">当前：{{ currentPaidFile.file_name }}</span>
            <input ref="paidPdfInput" type="file" accept=".pdf" style="display: none" @change="onPickPaidPdf" />
          </div>
          <div class="el-upload__tip">付费版 PDF 文件（高清版），且不超过 50MB</div>
        </el-form-item>

        <el-form-item label="封面图片">
          <div class="file-picker">
            <el-button size="small" @click="coverInput?.click()">
              {{ editingId ? '替换封面' : '选择封面' }}
            </el-button>
            <span v-if="coverFile" class="file-name">{{ coverFile.name }}</span>
            <input ref="coverInput" type="file" accept="image/jpeg,image/png,image/webp" style="display: none" @change="onPickCover" />
          </div>
          <div class="el-upload__tip">支持 JPG/PNG 格式，建议尺寸 800x800，且不超过 10MB</div>
          <el-image
            v-if="coverPreview"
            :src="coverPreview"
            fit="cover"
            style="width: 100px; height: 100px; margin-top: 8px; border-radius: 4px"
          />
        </el-form-item>

        <el-form-item label="音频文件">
          <div class="file-picker">
            <el-button size="small" @click="audioInput?.click()">
              {{ editingId ? '替换音频' : '选择音频' }}
            </el-button>
            <span v-if="audioFile" class="file-name">{{ audioFile.name }}</span>
            <input ref="audioInput" type="file" accept="audio/*" style="display: none" @change="onPickAudio" />
          </div>
          <div class="el-upload__tip">支持 MP3/WAV/FLAC/M4A 格式，且不超过 200MB</div>
          <audio
            v-if="!audioFile && currentAudioUrl"
            :src="currentAudioUrl"
            controls
            style="max-width: 100%; margin-top: 8px"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="saving" @click="handleSave">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { sheetMusicApi } from '@/api/sheet-music'
import { uploadSheetPDF, uploadAudio, uploadImage } from '@/api/file'
import type { SheetMusic, CreateSheetReq, Tag, SheetFile } from '@/types/sheet-music'

const sheets = ref<SheetMusic[]>([])
const allTags = ref<Tag[]>([])
const loading = ref(false)
const keyword = ref('')
const page = ref(1)
const pageSize = ref(20)
const total = ref(0)

const dialogVisible = ref(false)
const detailLoading = ref(false)
const editingId = ref<number | null>(null)
const saving = ref(false)
const formRef = ref()

const freePdfInput = ref<HTMLInputElement>()
const paidPdfInput = ref<HTMLInputElement>()
const coverInput = ref<HTMLInputElement>()
const audioInput = ref<HTMLInputElement>()
const freePdfFile = ref<File | null>(null)
const paidPdfFile = ref<File | null>(null)
const coverFile = ref<File | null>(null)
const audioFile = ref<File | null>(null)

const currentFreeFile = ref<SheetFile | null>(null)
const currentPaidFile = ref<SheetFile | null>(null)
const currentAudioUrl = ref('')

const emptyForm = (): CreateSheetReq => ({
  title: '',
  track_name: '',
  composer: '',
  arranger: '',
  transcriber: '',
  description: '',
  cover_url: '',
  external_link: '',
  price: 0,
  download_points: 0,
  tags: [],
})
const form = ref<CreateSheetReq>(emptyForm())

const staticBaseUrl = import.meta.env.VITE_API_RESOURCE_URL || ''
const toStaticUrl = (path: string) => {
  if (!path || path.startsWith('http') || path.startsWith('blob:')) return path
  return staticBaseUrl.replace(/\/$/, '') + path
}

const localCoverPreview = ref('')
const coverPreview = computed(() => {
  if (localCoverPreview.value) return localCoverPreview.value
  return form.value.cover_url ? toStaticUrl(form.value.cover_url) : ''
})

const formatTime = (t: string) => (t ? t.slice(0, 19).replace('T', ' ') : '')

const fetchTags = async () => {
  try {
    const res = await sheetMusicApi.getTags()
    allTags.value = res.result || []
  } catch {
    // 错误提示已由 request 拦截器统一处理
  }
}

const fetchSheets = async () => {
  loading.value = true
  try {
    const res = await sheetMusicApi.getList({
      page: page.value,
      page_size: pageSize.value,
      keyword: keyword.value || undefined,
    })
    if (res.result) {
      sheets.value = res.result.list || []
      total.value = res.result.total
    }
  } catch {
    // 错误提示已由 request 拦截器统一处理
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  page.value = 1
  fetchSheets()
}

const handlePageChange = (val: number) => {
  page.value = val
  fetchSheets()
}

const resetFiles = () => {
  freePdfFile.value = null
  paidPdfFile.value = null
  coverFile.value = null
  audioFile.value = null
  localCoverPreview.value = ''
  currentFreeFile.value = null
  currentPaidFile.value = null
  currentAudioUrl.value = ''
  for (const input of [freePdfInput.value, paidPdfInput.value, coverInput.value, audioInput.value]) {
    if (input) input.value = ''
  }
}

const openCreate = () => {
  editingId.value = null
  form.value = emptyForm()
  resetFiles()
  dialogVisible.value = true
}

const openEdit = async (row: SheetMusic) => {
  editingId.value = row.id
  resetFiles()
  form.value = {
    title: row.title,
    track_name: row.track_name || '',
    composer: row.composer || '',
    arranger: row.arranger || '',
    transcriber: row.transcriber || '',
    description: row.description || '',
    cover_url: row.cover_url || '',
    external_link: row.external_link || '',
    price: row.price || 0,
    download_points: row.download_points || 0,
    tags: [],
  }
  dialogVisible.value = true

  detailLoading.value = true
  try {
    const res = await sheetMusicApi.getDetail(row.id)
    const detail = res.result
    if (detail) {
      form.value = { ...form.value, tags: (detail.tags || []).map((t) => t.name) }
      const files = detail.files || []
      currentFreeFile.value = files.find((f) => f.file_type !== 'paid') || null
      currentPaidFile.value = files.find((f) => f.file_type === 'paid') || null
      const audio = detail.audio?.[0]
      currentAudioUrl.value = audio ? toStaticUrl(audio.hls_url || audio.original_url) : ''
    }
  } catch {
    // 错误提示已由 request 拦截器统一处理
  } finally {
    detailLoading.value = false
  }
}

const PDF_MAX_MB = 50
const IMAGE_MAX_MB = 10
const AUDIO_MAX_MB = 200

const pickFile = (e: Event, maxMB: number, tip: string): File | null => {
  const file = (e.target as HTMLInputElement).files?.[0]
  if (!file) return null
  if (file.size / 1024 / 1024 > maxMB) {
    ElMessage.error(`${tip}大小不能超过 ${maxMB}MB`)
    return null
  }
  return file
}

const onPickFreePdf = (e: Event) => {
  freePdfFile.value = pickFile(e, PDF_MAX_MB, '免费版 PDF')
}
const onPickPaidPdf = (e: Event) => {
  paidPdfFile.value = pickFile(e, PDF_MAX_MB, '付费版 PDF')
}
const onPickCover = (e: Event) => {
  const file = pickFile(e, IMAGE_MAX_MB, '封面图片')
  coverFile.value = file
  localCoverPreview.value = file ? URL.createObjectURL(file) : ''
}
const onPickAudio = (e: Event) => {
  audioFile.value = pickFile(e, AUDIO_MAX_MB, '音频')
}

const handleSave = async () => {
  if (!form.value.title) {
    ElMessage.warning('请填写标题')
    return
  }
  if (!editingId.value && !freePdfFile.value) {
    ElMessage.warning('请上传免费版乐谱 PDF')
    return
  }

  try {
    saving.value = true

    if (coverFile.value) {
      const res = await uploadImage(coverFile.value)
      if (res.result?.url) {
        form.value = { ...form.value, cover_url: res.result.url }
      }
    }

    let sheetId = editingId.value
    if (sheetId) {
      await sheetMusicApi.update(sheetId, form.value)
    } else {
      const res = await sheetMusicApi.create(form.value)
      sheetId = res.result?.id ?? null
      if (!sheetId) {
        ElMessage.error('创建失败')
        return
      }
    }

    if (freePdfFile.value) await uploadSheetPDF(sheetId, freePdfFile.value, 'free')
    if (paidPdfFile.value) await uploadSheetPDF(sheetId, paidPdfFile.value, 'paid')
    if (audioFile.value) await uploadAudio(sheetId, audioFile.value)

    ElMessage.success('保存成功')
    dialogVisible.value = false
    fetchSheets()
    fetchTags()
  } catch {
    // 错误提示已由 request 拦截器统一处理
  } finally {
    saving.value = false
  }
}

const handleDelete = async (row: SheetMusic) => {
  try {
    await ElMessageBox.confirm(`确定要删除乐谱「${row.title}」吗？`, '警告', { type: 'warning' })
    await sheetMusicApi.delete(row.id)
    ElMessage.success('删除成功')
    fetchSheets()
  } catch {
    // 取消或失败
  }
}

onMounted(() => {
  fetchSheets()
  fetchTags()
})
</script>

<style scoped>
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.header-actions {
  display: flex;
  align-items: center;
}
.pagination {
  display: flex;
  justify-content: center;
  margin-top: 20px;
}
.file-picker {
  display: flex;
  align-items: center;
  gap: 12px;
}
.file-name {
  color: #409eff;
  font-size: 13px;
}
.file-current {
  color: #909399;
  font-size: 13px;
}
.el-upload__tip {
  color: #909399;
  font-size: 12px;
  line-height: 1.5;
}
</style>
