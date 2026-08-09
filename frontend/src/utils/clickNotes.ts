/**
 * 鼠标点击音符特效：在点击处生成一个音符，向上漂浮并渐隐。
 * 动画样式定义在 global.scss 的 .click-note / @keyframes note-float。
 */

const NOTES = ['\u266A', '\u266B', '\u266C', '\u2669'] // ♪ ♫ ♬ ♩
const TREBLE_CLEF = '\u{1D11E}' // 𝄞，低概率彩蛋

// 暖纸色系，与全站主题一致
const COLORS = ['#C88D5E', '#A06B42', '#B87A4E', '#D4A351', '#8B7E6E']

const MAX_ACTIVE_NOTES = 24

let activeCount = 0
let initialized = false

function spawnNote(x: number, y: number) {
  if (activeCount >= MAX_ACTIVE_NOTES) return

  const note = document.createElement('span')
  note.className = 'click-note'
  note.textContent =
    Math.random() < 0.06 ? TREBLE_CLEF : NOTES[Math.floor(Math.random() * NOTES.length)]

  const size = 16 + Math.random() * 10
  const drift = (Math.random() - 0.5) * 48
  const rotate = (Math.random() - 0.5) * 30

  note.style.left = `${x}px`
  note.style.top = `${y - 8}px`
  note.style.fontSize = `${size}px`
  note.style.color = COLORS[Math.floor(Math.random() * COLORS.length)]
  note.style.setProperty('--note-drift', `${drift}px`)
  note.style.setProperty('--note-rotate', `${rotate}deg`)

  activeCount++
  let cleaned = false
  const cleanup = () => {
    if (cleaned) return
    cleaned = true
    note.remove()
    activeCount--
  }
  note.addEventListener('animationend', cleanup, { once: true })
  // 动画事件极端情况下不触发（如标签页隐藏）时的兜底回收
  setTimeout(cleanup, 1500)

  document.body.appendChild(note)
}

export function initClickNotes() {
  if (initialized) return
  initialized = true

  document.addEventListener('click', (e: MouseEvent) => {
    // 键盘触发的合成点击没有坐标，跳过
    if (e.clientX === 0 && e.clientY === 0) return
    spawnNote(e.clientX, e.clientY)
  })
}
