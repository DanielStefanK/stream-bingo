import { ref, watch } from 'vue'

const useDarkMode = () => {
  const isDarkMode = ref(false)
  const toggleDarkMode = () => {
    isDarkMode.value = !isDarkMode.value
    if (isDarkMode.value) {
      document.documentElement.classList.add('dark-mode')
    } else {
      document.documentElement.classList.remove('dark-mode')
    }
  }

  //check for dark mode preference
  const darkModeMediaQuery = window.matchMedia('(prefers-color-scheme: dark)')
  isDarkMode.value = darkModeMediaQuery.matches

  darkModeMediaQuery.addEventListener('change', (e) => {
    isDarkMode.value = e.matches
  })

  watch(
    isDarkMode,
    (value) => {
      if (value) {
        document.documentElement.classList.add('dark-mode')
      } else {
        document.documentElement.classList.remove('dark-mode')
      }
    },
    { immediate: true }
  )

  return { isDarkMode, toggleDarkMode }
}

export default useDarkMode
