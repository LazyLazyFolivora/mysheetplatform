<template>
  <div class="sheet-music-detail">
    <div v-if="loading" class="loading-container">
      <el-skeleton :rows="10" animated />
    </div>

    <div v-else-if="sheet" class="detail-layout">
      <!-- 左侧主要内容 -->
      <div class="main-content">
        <!-- 乐谱预览：逐页图片轮播 -->
        <div class="preview-container">
          <template v-if="sheetImages.length">
            <div class="zoom-controls">
              <el-button-group>
                <el-button @click="decreaseZoom" :disabled="zoom <= ZOOM_MIN">
                  <el-icon><ZoomOut /></el-icon>
                </el-button>
                <el-button @click="resetZoom">{{ zoom }}%</el-button>
                <el-button @click="increaseZoom" :disabled="zoom >= ZOOM_MAX">
                  <el-icon><ZoomIn /></el-icon>
                </el-button>
              </el-button-group>
            </div>

            <el-carousel :autoplay="false" trigger="click" indicator-position="outside">
              <el-carousel-item v-for="(image, index) in sheetImages" :key="index">
                <div class="sheet-page">
                  <img
                    :src="image"
                    class="sheet-page-image"
                    @contextmenu.prevent
                    :alt="`乐谱第 ${index + 1} 页`"
                    :style="{ width: `${zoom}%` }"
                  />
                </div>
              </el-carousel-item>
            </el-carousel>
          </template>
          <!-- 图片未生成时回退到 PDF 内嵌预览 -->
          <iframe v-else-if="freePdfUrl" :src="freePdfUrl" class="pdf-frame" title="乐谱预览"></iframe>
          <el-empty v-else description="暂无乐谱预览" />
        </div>

        <!-- 试听部分 -->
        <div class="audio-section">
          <h3>试听</h3>
          <audio-player v-if="audioUrl" :audio-url="audioUrl" />
          <el-empty v-else description="暂无试听音频" />
        </div>

        <!-- 下载/购买按钮 -->
        <div class="download-section">
          <el-button
            type="primary"
            size="large"
            @click="handleDownload"
            :loading="downloading"
            class="download-btn"
          >
            下载乐谱 (免费版)
          </el-button>

          <el-button
            v-if="sheet.price > 0"
            type="success"
            size="large"
            @click="handlePurchase"
            :loading="purchasing"
            class="purchase-btn"
          >
            购买高清版 (¥{{ sheet.price.toFixed(2) }})
          </el-button>
        </div>

        <!-- 外链部分 -->
        <div v-if="sheet.external_link" class="external-link-section">
          <h3>外部链接</h3>
          <iframe :src="sheet.external_link" frameborder="0" class="external-frame"></iframe>
        </div>
      </div>

      <!-- 右侧信息卡片 -->
      <div class="info-card">
        <el-card>
          <template #header>
            <h2>{{ sheet.title }}</h2>
          </template>
          <div class="info-content">
            <div class="info-item">
              <span class="label">作曲</span>
              <span>{{ sheet.composer || '未知' }}</span>
            </div>
            <div class="info-item">
              <span class="label">编曲</span>
              <span>{{ sheet.arranger || '未知' }}</span>
            </div>
            <div class="info-item">
              <span class="label">扒谱</span>
              <span>{{ sheet.transcriber || '未知' }}</span>
            </div>
            <div class="info-item">
              <span class="label">下载积分</span>
              <span>{{ sheet.download_points || 0 }}</span>
            </div>
            <div class="info-item">
              <span class="label">价格</span>
              <span>¥{{ (sheet.price || 0).toFixed(2) }}</span>
            </div>
            <div class="info-item">
              <span class="label">下载次数</span>
              <span>{{ freeFile?.download_count || 0 }}</span>
            </div>
            <div class="info-item">
              <span class="label">浏览次数</span>
              <span>{{ sheet.view_count || 0 }}</span>
            </div>
            <div class="info-item" v-if="sheet.tags?.length">
              <span class="label">标签</span>
              <div class="tags">
                <el-tag v-for="tag in sheet.tags" :key="tag.id" size="small" class="tag">
                  {{ tag.name }}
                </el-tag>
              </div>
            </div>
            <div class="info-item description" v-if="sheet.description">
              <span class="label">描述</span>
              <p>{{ sheet.description }}</p>
            </div>
          </div>
        </el-card>
      </div>
    </div>

    <el-empty v-else description="未找到乐谱信息" />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ZoomIn, ZoomOut } from '@element-plus/icons-vue'
import { sheetMusicApi } from '@/api/sheet-music'
import type { SheetDetail } from '@/types/sheet-music'
import AudioPlayer from '@/components/AudioPlayer.vue'
import { useUserStore } from '@/stores/user'

const ZOOM_MIN = 20
const ZOOM_MAX = 200
const ZOOM_STEP = 10
const ZOOM_DEFAULT = 80

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()
const sheet = ref<SheetDetail | null>(null)
const loading = ref(true)
const downloading = ref(false)
const purchasing = ref(false)
const zoom = ref(ZOOM_DEFAULT)

const staticBaseUrl = import.meta.env.VITE_API_RESOURCE_URL || ''
function getStaticUrl(path: string) {
  if (!path) return ''
  if (path.startsWith('http')) return path
  return staticBaseUrl.replace(/\/$/, '') + path
}

const freeFile = computed(() => sheet.value?.files?.find((f) => f.file_type !== 'paid') || null)

const freePdfUrl = computed(() => (freeFile.value ? getStaticUrl(freeFile.value.file_path) : ''))

const sheetImages = computed(() => {
  const file = freeFile.value
  if (!file?.image_folder || !file.page_count) return []
  return Array.from({ length: file.page_count }, (_, i) =>
    getStaticUrl(`${file.image_folder}${i + 1}.png`)
  )
})

const increaseZoom = () => {
  zoom.value = Math.min(zoom.value + ZOOM_STEP, ZOOM_MAX)
}

const decreaseZoom = () => {
  zoom.value = Math.max(zoom.value - ZOOM_STEP, ZOOM_MIN)
}

const resetZoom = () => {
  zoom.value = ZOOM_DEFAULT
}

const audioUrl = computed(() => {
  const audio = sheet.value?.audio?.[0]
  if (!audio) return ''
  return getStaticUrl(audio.hls_url || audio.original_url || '')
})

const handleDownload = async () => {
  if (!userStore.isLoggedIn) {
    ElMessage.warning('请先登录后再下载')
    router.push('/login')
    return
  }
  if (!freeFile.value) {
    ElMessage.warning('暂无可下载的PDF文件')
    return
  }
  const requiredPoints = sheet.value?.download_points || 0
  try {
    await ElMessageBox.confirm(
      `确认下载该乐谱吗？\n所需积分：${requiredPoints} 积分`,
      '下载确认',
      { confirmButtonText: '确认下载', cancelButtonText: '取消', type: 'warning' }
    )
    downloading.value = true
    const token = localStorage.getItem('token') || ''
    const baseUrl = import.meta.env.VITE_API_BASE_URL || '/api'
    const downloadUrl = `${baseUrl}/sheet-music/${sheet.value!.id}/download?token=${encodeURIComponent(token)}`
    window.open(downloadUrl, '_blank')
    ElMessage.success('下载已开始，请在浏览器下载列表中查看')
    // 刷新用户积分（下载可能扣了积分）
    setTimeout(() => userStore.fetchProfile?.(), 1000)
  } catch (err: any) {
    if (err === 'cancel') return
  } finally {
    downloading.value = false
  }
}

const handlePurchase = async () => {
  if (!userStore.isLoggedIn) {
    ElMessage.warning('请先登录后再购买')
    router.push('/login')
    return
  }
  if (!sheet.value) return
  try {
    await ElMessageBox.confirm(
      `确认购买该乐谱吗？\n价格：¥${sheet.value.price.toFixed(2)}`,
      '购买确认',
      { confirmButtonText: '确认购买', cancelButtonText: '取消', type: 'warning' }
    )
    purchasing.value = true
    const res = await sheetMusicApi.createOrder(sheet.value.id)
    if (res.result) {
      ElMessage.success(`订单已创建（单号 ${res.result.order_no}），正在跳转支付宝支付...`)
      const token = localStorage.getItem('token') || ''
      const payUrl = import.meta.env.VITE_API_BASE_URL
        ? `${import.meta.env.VITE_API_BASE_URL}/orders/${res.result.order_no}/pay?token=${encodeURIComponent(token)}`
        : `/api/orders/${res.result.order_no}/pay?token=${encodeURIComponent(token)}`
      window.open(payUrl, '_blank')
    }
  } catch {
    // 取消或失败（失败提示已由拦截器处理）
  } finally {
    purchasing.value = false
  }
}

onMounted(async () => {
  try {
    const res = await sheetMusicApi.getDetail(Number(route.params.id))
    if (res.result) {
      sheet.value = res.result
    }
  } catch {
    // 错误提示已由 request 拦截器统一处理
  } finally {
    loading.value = false
  }
})
</script>

<style scoped>
.loading-container {
  padding: 20px;
}

.detail-layout {
  display: flex;
  gap: 24px;
  margin: 20px auto;
  max-width: 1600px;
  padding: 0 20px;
}

.main-content {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.preview-container {
  position: relative;
  width: 100%;
  background-color: #f5f7fa;
  border-radius: 4px;
  overflow: hidden;
}

.zoom-controls {
  position: absolute;
  top: 20px;
  right: 20px;
  z-index: 10;
}

.sheet-page {
  display: flex;
  justify-content: center;
  align-items: center;
  background-color: #f5f7fa;
  min-height: 1500px;
}

.sheet-page-image {
  transition: width 0.3s ease;
  object-fit: contain;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
  background-color: white;
  border-radius: 4px;
}

:deep(.el-carousel__container) {
  height: 1500px;
}

:deep(.el-carousel__item) {
  display: flex;
  justify-content: center;
  align-items: center;
  overflow-y: auto;
}

.pdf-frame {
  width: 100%;
  height: 1000px;
  border: none;
  display: block;
}

.audio-section,
.external-link-section {
  background: white;
  padding: 20px;
  border-radius: 4px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
}

.audio-section h3,
.external-link-section h3 {
  margin: 0 0 12px;
}

.download-section {
  display: flex;
  justify-content: center;
  padding: 20px 0;
  gap: 10px;
}

.download-btn,
.purchase-btn {
  min-width: 200px;
}

.external-frame {
  width: 100%;
  height: 600px;
  margin-top: 10px;
}

.info-card {
  width: 360px;
  flex-shrink: 0;
}

.info-card h2 {
  margin: 0;
  font-size: 20px;
}

.info-content {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.info-item {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.label {
  font-weight: bold;
  color: #606266;
}

.tags {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.description p {
  margin: 0;
  line-height: 1.6;
  color: #606266;
  white-space: pre-wrap;
}

@media (max-width: 900px) {
  .detail-layout {
    flex-direction: column;
  }
  .info-card {
    width: 100%;
  }
  .pdf-frame {
    height: 600px;
  }
}
</style>
