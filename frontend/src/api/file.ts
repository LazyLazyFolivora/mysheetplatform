import request from '@/utils/request'
import type { ApiResponse } from '@/types/api'

// Go backend file upload endpoints
export function uploadSheetPDF(sheetMusicId: number, file: File, fileType: 'free' | 'paid' = 'free') {
  const formData = new FormData()
  formData.append('sheet_music_id', String(sheetMusicId))
  formData.append('file_type', fileType)
  formData.append('file', file)
  return request.post<ApiResponse<any>>('/admin/sheet-music/upload-pdf', formData, {
    headers: { 'Content-Type': 'multipart/form-data' },
  })
}

export function uploadAudio(sheetMusicId: number, file: File) {
  const formData = new FormData()
  formData.append('sheet_music_id', String(sheetMusicId))
  formData.append('file', file)
  return request.post<ApiResponse<any>>('/admin/sheet-music/upload-audio', formData, {
    headers: { 'Content-Type': 'multipart/form-data' },
  })
}

export function uploadImage(file: File) {
  const formData = new FormData()
  formData.append('file', file)
  return request.post<ApiResponse<{ url: string }>>('/admin/upload/image', formData, {
    headers: { 'Content-Type': 'multipart/form-data' },
  })
}
