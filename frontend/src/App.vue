<script setup lang="ts">
import useDarkMode from './composables/useDarkMode'
import Login from './pages/LoginPage.vue'

useDarkMode()
</script>

<template>
  <Login />
</template>

<style scoped></style>
