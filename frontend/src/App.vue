<script setup lang="ts">
import useDarkMode from './composables/useDarkMode'
import Login from './pages/Login.vue'

useDarkMode()
</script>

<template>
  <Login />
</template>

<style scoped></style>
