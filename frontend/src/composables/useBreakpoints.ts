import { ref } from 'vue'

const useBreakpoints = () => {
  const mobileQuery = window.matchMedia('(max-width: 640px)')
  const tabletQuery = window.matchMedia('(max-width: 768px)')
  const desktopQuery = window.matchMedia('(min-width: 1024px)')
  const isMobile = ref(mobileQuery.matches)
  const isTablet = ref(tabletQuery.matches)
  const isDesktop = ref(desktopQuery.matches)

  //on resize
  const handleResize = () => {
    isMobile.value = mobileQuery.matches
    isTablet.value = tabletQuery.matches
    isDesktop.value = desktopQuery.matches
  }

  //TODO: memory leak
  window.addEventListener('resize', handleResize)

  return {
    isMobile,
    isTablet,
    isDesktop,
  }
}

export default useBreakpoints
