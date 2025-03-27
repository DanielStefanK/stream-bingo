<script lang="ts" setup>
import ProgressSpinner from 'primevue/progressspinner'
import { loginOAuth } from '../services/auth'
import { RouterLink, useRoute, useRouter } from 'vue-router'
import { onMounted, ref } from 'vue'
import { useUserStore } from '../store/userStore'
import { Button, Message } from 'primevue'

const route = useRoute()
const router = useRouter()
const error = ref('')
const userStore = useUserStore()
const isLoading = ref(true)

const oauthCallback = async () => {
  try {
    const resp = await loginOAuth(
      route.params.provider as string,
      route.query.code as string
    )
    userStore.$state.user = resp.user
    userStore.$state.token = resp.token
    router.push('/')
  } catch (err) {
    console.error(err)
    error.value = 'An error occurred while authenticating'
  } finally {
    isLoading.value = false
  }
}

onMounted(() => {
  oauthCallback()
})
</script>
<template>
  <div class="center">
    <div v-if="isLoading" class="loading">
      <p>Authenticating...</p>
      <ProgressSpinner />
      <p>you will be redirected shortly</p>
    </div>
    <div v-else class="loading">
      <Message severity="error">
        {{ error }}
      </Message>
      <Button :as="RouterLink" to="/login" class="unstyled back-btn">
        <i class="pi pi-arrow-left"></i>
        back to login</Button
      >
    </div>
  </div>
</template>

<style scoped>
.center {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100vh;
}
.loading {
  display: flex;
  flex-direction: column;
  align-items: center;
}
.back-btn {
  margin-top: 1rem;
}
</style>
