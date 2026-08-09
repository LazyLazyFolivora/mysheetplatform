<template>
  <div class="audio-player">
    <audio ref="audioRef" controls style="width: 100%;" @timeupdate="handleTimeUpdate"></audio>
  </div>
</template>

<script setup lang="ts">
import { ref, watch, onMounted, onUnmounted } from 'vue'
import Hls from 'hls.js'

const props = defineProps<{
  audioUrl: string
}>()

const emit = defineEmits<{
  (e: 'timeupdate', currentTime: number): void
}>()

const handleTimeUpdate = () => {
  emit('timeupdate', audioRef.value?.currentTime ?? 0)
}

const audioRef = ref<HTMLAudioElement | null>(null)
let hls: Hls | null = null

const baseUrl = import.meta.env.VITE_API_RESOURCE_URL || ''
const getFullUrl = (url: string) => {
  if (url.startsWith('http')) {
    return url
  }
  return `${baseUrl.replace(/\/$/, '')}${url}`
}

const withToken = (url: string, token: string | null) => {
  if (!token) return url
  const separator = url.includes('?') ? '&' : '?'
  return `${url}${separator}token=${token}`
}

const initPlayer = () => {
  if (hls) {
    hls.destroy()
    hls = null
  }

  const audio = audioRef.value
  if (!audio || !props.audioUrl) return

  const fullUrl = getFullUrl(props.audioUrl)
  const token = localStorage.getItem('token')

  if (fullUrl.endsWith('.m3u8')) {
    if (Hls.isSupported()) {
      const hlsConfig: any = {}
      if (token) {
        hlsConfig.xhrSetup = (xhr: XMLHttpRequest, url: string) => {
          xhr.setRequestHeader('Authorization', `Bearer ${token}`)
          if (!url.startsWith('http')) {
            xhr.open('GET', getFullUrl(url), true)
          }
        }
      }

      hls = new Hls(hlsConfig)
      hls.loadSource(fullUrl)
      hls.attachMedia(audio)
      hls.on(Hls.Events.ERROR, (_event, data) => {
        if (data.type === 'networkError' && data.details === 'manifestLoadError') {
          setTimeout(() => {
            hls?.loadSource(fullUrl)
          }, 1000)
        }
      })
    } else if (audio.canPlayType('application/vnd.apple.mpegurl')) {
      // Safari 原生 HLS 无法设置请求头，只能通过 query 传 token
      audio.src = withToken(fullUrl, token)
    }
  } else {
    audio.src = withToken(fullUrl, token)
  }
}

watch(() => props.audioUrl, (newUrl) => {
  if (newUrl) {
    initPlayer()
  }
}, { immediate: true })

onMounted(() => {
  initPlayer()
})

onUnmounted(() => {
  if (hls) {
    hls.destroy()
    hls = null
  }
})
</script>

<style scoped>
.audio-player {
  width: 100%;
}
</style>
