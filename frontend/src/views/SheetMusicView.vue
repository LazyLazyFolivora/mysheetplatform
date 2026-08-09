<template>
  <div class="sheet-music-container">
    <div class="search-section">
      <el-input
        v-model="searchQuery"
        placeholder="搜索乐谱名称或作曲家"
        class="search-input"
        clearable
        @keyup.enter="handleSearch"
      >
        <template #prefix>
          <el-icon><search /></el-icon>
        </template>
      </el-input>
      <el-select
        v-model="selectedTagIds"
        multiple
        filterable
        placeholder="按标签筛选"
        class="tag-select"
        clearable
      >
        <el-option v-for="tag in availableTags" :key="tag.id" :label="tag.name" :value="tag.id" />
      </el-select>
      <el-button type="primary" @click="handleSearch">搜索</el-button>
    </div>

    <div v-if="loading" class="loading-section">
      <el-skeleton :rows="8" animated />
    </div>

    <div v-else-if="sheetList.length > 0" class="sheet-music-list">
      <el-row :gutter="20">
        <el-col
          v-for="sheet in sheetList"
          :key="sheet.id"
          :xs="24"
          :sm="12"
          :md="8"
          :lg="6"
        >
          <el-card
            class="sheet-music-card"
            :body-style="{ padding: '0px' }"
            @click="goToDetail(sheet.id)"
          >
            <img :src="getCoverImage(sheet)" class="cover-image" @error="onImageError" />
            <div class="card-content">
              <h3 class="sheet-title">{{ sheet.title }}</h3>
              <p class="sheet-detail"><strong>作曲:</strong> {{ sheet.composer || 'N/A' }}</p>
              <div class="sheet-tags" v-if="sheet.tags && sheet.tags.length">
                <el-tag v-for="tag in sheet.tags" :key="tag.id" size="small" class="sheet-tag">
                  {{ tag.name }}
                </el-tag>
              </div>
              <div class="card-footer">
                <span><el-icon><Coin /></el-icon> {{ sheet.download_points || 0 }} 积分</span>
              </div>
            </div>
          </el-card>
        </el-col>
      </el-row>
    </div>
    <el-empty v-else description="暂无乐谱数据"></el-empty>

    <div class="pagination">
      <el-pagination
        :current-page="currentPage"
        :page-size="pageSize"
        :page-sizes="[8, 16, 24, 32]"
        :total="total"
        layout="total, sizes, prev, pager, next, jumper"
        @size-change="handleSizeChange"
        @current-change="handleCurrentChange"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { Search, Coin } from '@element-plus/icons-vue'
import { sheetMusicApi } from '@/api/sheet-music'
import type { SheetListItem, Tag } from '@/types/sheet-music'
import placeholderImage from '@/assets/logo.png'

const router = useRouter()
const searchQuery = ref('')
const currentPage = ref(1)
const pageSize = ref(8)
const sheetList = ref<SheetListItem[]>([])
const total = ref(0)
const loading = ref(false)
const availableTags = ref<Tag[]>([])
const selectedTagIds = ref<number[]>([])

const staticBaseUrl = import.meta.env.VITE_API_RESOURCE_URL || ''
function getStaticUrl(path: string) {
  if (!path) return ''
  if (path.startsWith('http')) return path
  return staticBaseUrl.replace(/\/$/, '') + path
}

const getCoverImage = (sheet: SheetListItem) => {
  if (sheet.cover_url) {
    return getStaticUrl(sheet.cover_url)
  }
  return placeholderImage
}

const onImageError = (e: Event) => {
  const target = e.target as HTMLImageElement
  target.src = placeholderImage
}

const handleSearch = () => {
  currentPage.value = 1
  fetchSheetMusic()
}

const handleSizeChange = (val: number) => {
  pageSize.value = val
  fetchSheetMusic()
}

const handleCurrentChange = (val: number) => {
  currentPage.value = val
  fetchSheetMusic()
}

const goToDetail = (id: number) => {
  router.push(`/sheet-music/${id}`)
}

const fetchSheetMusic = async () => {
  loading.value = true
  try {
    const response = await sheetMusicApi.getList({
      page: currentPage.value,
      page_size: pageSize.value,
      keyword: searchQuery.value || undefined,
      tag_ids: selectedTagIds.value.length ? selectedTagIds.value : undefined,
    })
    if (response.result) {
      sheetList.value = response.result.list || []
      total.value = response.result.total
    }
  } catch {
    // 错误提示已由 request 拦截器统一处理
  } finally {
    loading.value = false
  }
}

const fetchTags = async () => {
  try {
    const res = await sheetMusicApi.getTags()
    availableTags.value = res.result || []
  } catch {
    // 错误提示已由 request 拦截器统一处理
  }
}

onMounted(() => {
  fetchSheetMusic()
  fetchTags()
})
</script>

<style scoped lang="scss">
@use "@/styles/variables" as *;

.sheet-music-container {
  max-width: $content-max-width;
  margin: 0 auto;
  padding: 20px;
}

.search-section {
  display: flex;
  gap: 12px;
  margin-bottom: $space-xl;
}

.search-input {
  width: 300px;
}

.tag-select {
  width: 300px;
}

.sheet-tags {
  margin: 8px 0;
  display: flex;
  flex-wrap: wrap;
  gap: 4px;
}

.loading-section {
  padding: 20px;
  background: $bg-card;
  border: 1px solid $border-light;
  border-radius: $radius-lg;
  margin-bottom: 20px;
}

.sheet-music-list {
  margin-bottom: 20px;
}

.sheet-music-card {
  margin-bottom: 20px;
  cursor: pointer;
  display: flex;
  flex-direction: column;
  height: 100%;
  border-color: $border-light;
  border-radius: $radius-xl;
  overflow: hidden;
  transition: transform $transition-base, box-shadow $transition-base, border-color $transition-base;

  &:hover {
    transform: translateY(-2px);
    box-shadow: $shadow-card-hover;
    border-color: $primary-light;

    .cover-image {
      transform: scale(1.03);
    }
  }
}

.cover-image {
  width: 100%;
  height: 200px;
  object-fit: cover;
  border-bottom: 1px solid $border-light;
  transition: transform 0.35s ease;
}

.card-content {
  padding: 16px;
  flex-grow: 1;
}

.sheet-title {
  margin: 0 0 8px;
  font-size: 1.1rem;
  font-weight: 600;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.sheet-detail {
  margin: 4px 0;
  color: $text-secondary;
  font-size: 0.9rem;
}

.card-footer {
  margin-top: 12px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  color: $text-secondary;
  font-size: 0.85rem;

  span {
    display: inline-flex;
    align-items: center;
    gap: 4px;
    color: $primary-dark;
  }
}

.pagination {
  display: flex;
  justify-content: center;
  margin-top: 20px;
}
</style>
