<script setup lang="ts">
import { Toast } from 'primevue'
import useDarkMode from './composables/useDarkMode'

useDarkMode()
</script>

<template>
  <Toast />
  <RouterView />
</template>

<style scoped></style>
