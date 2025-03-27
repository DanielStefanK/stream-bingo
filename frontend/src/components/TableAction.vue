<script setup lang="ts">
import Menu from 'primevue/menu'
import Button from 'primevue/button'
import { ref } from 'vue'
import type { MenuItem } from 'primevue/menuitem'

const menu = ref()

defineProps<{
  items: MenuItem[]
  loading?: boolean
}>()

const toggle = (e: MouseEvent) => {
  menu.value.toggle(e)
}
</script>

<template>
  <div class="card flex justify-center">
    <Button
      aria-haspopup="true"
      aria-controls="overlay_menu"
      type="button"
      :disabled="loading"
      :icon="loading ? 'pi pi-spin pi-spinner' : 'pi pi-ellipsis-v'"
      @click="toggle"
    >
    </Button>
    <Menu id="overlay_menu" ref="menu" :model="items" :popup="true" />
  </div>
</template>
