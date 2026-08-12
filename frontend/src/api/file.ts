import request from '@/utils/request'
import type { ApiResponse } from '@/types/api'

// 不要手动设置 Content-Type：浏览器/axios 会自动带 multipart boundary。
// 若写成 'multipart/form-data' 且无 boundary，后端 PostForm 解析失败，
// file_type 会掉进默认值 free，导致付费 PDF 也被存成免费版。

export function uploadSheetPDF(sheetMusicId: number, file: File, fileType: 'free' | 'paid' = 'free') {
  const formData = new FormData()
  formData.append('sheet_music_id', String(sheetMusicId))
  formData.append('file_type', fileType)
  formData.append('file', file)
  return request.post<ApiResponse<any>>('/admin/sheet-music/upload-pdf', formData)
}

export function uploadAudio(sheetMusicId: number, file: File) {
  const formData = new FormData()
  formData.append('sheet_music_id', String(sheetMusicId))
  formData.append('file', file)
  return request.post<ApiResponse<any>>('/admin/sheet-music/upload-audio', formData)
}

export function uploadImage(file: File) {
  const formData = new FormData()
  formData.append('file', file)
  return request.post<ApiResponse<{ url: string }>>('/admin/upload/image', formData)
}
