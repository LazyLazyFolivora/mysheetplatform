import request from '@/utils/request'
import type { SheetListParams, SheetListResp, SheetDetail, SheetMusic, CreateSheetReq, Tag } from '@/types/sheet-music'
import type { ApiResponse } from '@/types/api'

export const sheetMusicApi = {
  getList: (params: SheetListParams) => {
    const queryParams = new URLSearchParams()
    queryParams.append('page', String(params.page || 1))
    queryParams.append('page_size', String(params.page_size || 20))
    if (params.keyword) queryParams.append('keyword', params.keyword)
    if (params.tag_ids && params.tag_ids.length > 0) {
      params.tag_ids.forEach((id) => queryParams.append('tag_ids', String(id)))
    }
    return request.get<ApiResponse<SheetListResp>>(`/sheet-music?${queryParams.toString()}`)
  },

  getDetail: (id: number) => request.get<ApiResponse<SheetDetail>>(`/sheet-music/${id}`),

  getTags: () => request.get<ApiResponse<Tag[]>>('/tags'),

  create: (data: CreateSheetReq) => request.post<ApiResponse<SheetMusic>>('/admin/sheet-music', data),

  update: (id: number, data: CreateSheetReq) =>
    request.put<ApiResponse<SheetMusic>>(`/admin/sheet-music/${id}`, data),

  delete: (id: number) => request.delete<ApiResponse<void>>(`/admin/sheet-music/${id}`),

  createOrder: (sheetMusicId: number) =>
    request.post<ApiResponse<{ id: number; order_no: string; amount: number; status: string; pay_url: string }>>('/orders', {
      sheet_music_id: sheetMusicId,
    }),

  getOrderStatus: (orderNo: string) =>
    request.get<ApiResponse<string>>('/alipay/order/status', { params: { order_no: orderNo } }),

  download: (sheetMusicId: number) =>
    request.get(`/sheet-music/${sheetMusicId}/download`, { responseType: 'blob' }),
}
