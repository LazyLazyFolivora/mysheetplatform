import { createApp } from 'vue'
import { createPinia } from 'pinia'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import * as ElementPlusIconsVue from '@element-plus/icons-vue'

import App from './App.vue'
import router from './router'
import { useUserStore } from './stores/user'

import '@/styles/global.scss'
import '@/styles/element-overrides.scss'
import { initClickNotes } from '@/utils/clickNotes'

const app = createApp(App)

for (const [key, component] of Object.entries(ElementPlusIconsVue)) {
  app.component(key, component)
}

const pinia = createPinia()
app.use(pinia)
app.use(router)
app.use(ElementPlus)

// Restore user session on boot
const userStore = useUserStore()
// initUserInfo is called inside the store constructor automatically

initClickNotes()

app.mount('#app')
