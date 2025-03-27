<script lang="ts" setup>
import { ref } from 'vue'
import { useUserStore } from '../store/userStore'
import { RouterLink, useRoute } from 'vue-router'
import Divider from 'primevue/divider'
import Button from 'primevue/button'
import { Card } from 'primevue'
import useBreakpoints from '../composables/useBreakpoints'
import Drawer from 'primevue/drawer'

const userStore = useUserStore()
const routeMeta = useRoute().meta
const { isMobile } = useBreakpoints()

const visible = ref(false)

const items = ref([
  {
    id: 'dashboard',
    to: '/',
    label: 'Home',
    icon: 'pi pi-home',
  },

  ...(userStore.user?.admin
    ? [
        {
          id: 'admin-separator',
          separator: true,
        },
        {
          id: 'admin-users',
          to: '/admin/user/list',
          label: 'Show Users',
          icon: 'pi pi-users',
        },
      ]
    : []),
])
</script>

<template>
  <div class="layout">
    <div class="sideBar">
      <Card
        v-if="!isMobile"
        class="sideBar-card"
        :pt="{ body: { style: 'padding:0' } }"
      >
        <template #content>
          <div class="btn-list">
            <template v-for="item in items" :key="item.id">
              <template v-if="item.separator">
                <Divider />
              </template>
              <template v-else>
                <Button
                  v-tooltip.right="item.label"
                  :color="item.to === $route.path ? 'primary' : 'secondary'"
                  :as="RouterLink"
                  :icon="item.icon"
                  :to="item.to"
                  class="unstyled"
                >
                </Button>
              </template>
            </template>
          </div>
        </template>
      </Card>
      <Drawer v-model:visible="visible">
        <div class="btn-list">
          <template v-for="item in items" :key="item.id">
            <template v-if="item.separator">
              <Divider />
            </template>
            <template v-else>
              <Button
                v-tooltip.right="item.label"
                :color="item.to === $route.path ? 'primary' : 'secondary'"
                :as="RouterLink"
                :to="item.to"
                class="unstyled"
              >
                <i :class="item.icon"></i>
                <span>{{ item.label }}</span>
              </Button>
            </template>
          </template>
        </div>
      </Drawer>
    </div>
    <div class="title">
      <div class="title-start">
        <Button v-if="isMobile" icon="pi pi-bars" @click="visible = true" />
        <div>
          <span class="sub-page-title">{{ routeMeta.subTitle }}</span>
          <h1 class="page-title">{{ routeMeta.title }}</h1>
        </div>
      </div>
      <div>
        <slot name="actions"></slot>
      </div>
    </div>
    <div class="mainContent"><slot /></div>
  </div>
</template>

<style scoped>
.layout {
  display: grid;
  grid-template-columns: 3.5rem 1fr;
  grid-template-rows: 5rem 1fr;
  margin: 1rem;
  height: 100%;
  gap: 1rem;
}

@media (max-width: 640px) {
  .layout {
    grid-template-columns: 0rem 1fr;

    grid-template-rows: 5rem 1fr;
  }
}
.sideBar {
  grid-column: 1 / 2;
  grid-row: 1 / 3;
  height: 100%;
}
.title-start {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.sideBar-card {
  max-width: 3.5rem;
  height: 100%;
}

.title {
  grid-column: 2 / 3;
  grid-row: 1 / 2;
  display: flex;
  justify-content: space-between;
  font-size: 1.5rem;
  font-weight: bold;
}

.page-title {
  font-size: 1.5rem;
  font-weight: bold;
  margin: 0;
}

.sub-page-title {
  font-size: 1rem;
  font-weight: 100;
  color: var(--p-surface-400);
}

.mainContent {
  grid-column: 2 / 3;
  grid-row: 2 / 3;
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.btn-list {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  justify-items: center;
  gap: 0.5rem;
  margin: 0.5rem;
}
</style>
