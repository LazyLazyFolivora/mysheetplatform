export interface Tag {
  id: number
  name: string
  type?: string
}

/** 曲谱同步点：播放到 time 秒时翻到第 page 页 */
export interface PageSyncPoint {
  time: number
  page: number
}

export interface SheetMusic {
  id: number
  title: string
  track_name?: string
  composer?: string
  arranger?: string
  transcriber?: string
  description?: string
  cover_url?: string
  external_link?: string
  price: number
  download_points: number
  status: number
  view_count: number
  user_id: number
  created_at: string
  updated_at: string
  /** 曲谱同步点（可选），为空表示不启用播放翻页同步 */
  page_sync?: PageSyncPoint[]
}

export interface SheetFile {
  id: number
  sheet_music_id: number
  file_path: string
  file_name: string
  file_type: 'free' | 'paid'
  image_folder: string
  page_count: number
  file_size: number
  points: number
  download_count: number
  created_at: string
}

export interface AudioFile {
  id: number
  sheet_music_id: number
  original_url: string
  hls_url: string
  quality: string
  duration: number
  file_size: number
}

export interface SheetListParams {
  page?: number
  page_size?: number
  keyword?: string
  tag_ids?: number[]
}

export interface SheetListItem extends SheetMusic {
  tags: Tag[]
}

export interface SheetListResp {
  list: SheetListItem[]
  total: number
  page: number
  size: number
}

export interface SheetDetail extends SheetMusic {
  tags: Tag[]
  files?: SheetFile[]
  audio?: AudioFile[]
  like_count: number
  is_liked: boolean
  is_purchased: boolean
  has_free_file: boolean
  has_paid_file: boolean
}

export interface CreateSheetReq {
  title: string
  track_name?: string
  composer?: string
  arranger?: string
  transcriber?: string
  description?: string
  cover_url?: string
  external_link?: string
  price?: number
  download_points?: number
  tags?: string[]
  page_sync?: PageSyncPoint[]
}
