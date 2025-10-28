// Keyboard navigation composable
// Handles arrow key navigation for date changes

import { onMounted, onUnmounted } from 'vue'

export function useKeyboardNavigation(
  onPrevious: () => void,
  onNext: () => void,
  enabled = true
) {
  const handleKeydown = (event: KeyboardEvent) => {
    if (!enabled) return

    // Ignore if user is typing in an input field
    const target = event.target as HTMLElement
    if (target.tagName === 'INPUT' || target.tagName === 'TEXTAREA') {
      return
    }

    switch (event.key) {
      case 'ArrowLeft':
        event.preventDefault()
        onPrevious()
        break
      case 'ArrowRight':
        event.preventDefault()
        onNext()
        break
    }
  }

  onMounted(() => {
    window.addEventListener('keydown', handleKeydown)
  })

  onUnmounted(() => {
    window.removeEventListener('keydown', handleKeydown)
  })

  return {
    handleKeydown
  }
}
