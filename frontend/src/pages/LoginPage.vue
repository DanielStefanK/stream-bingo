<script lang="ts" setup>
import { ref } from 'vue'
import LoginCard from '../components/LoginCard.vue'
import { loginLocal, register } from '../services/auth'
import { ApiError } from '../services/api'
import { useUserStore } from '../store/userStore'
import { useRouter } from 'vue-router'
import { useToast } from 'primevue'

const email = ref('')
const password = ref('')
const name = ref('')
const signupMode = ref(false)
const userStore = useUserStore()
const router = useRouter()
const toast = useToast()

const login = async () => {
  try {
    const resp = await loginLocal(email.value, password.value)
    userStore.$state.user = resp.user
    userStore.$state.token = resp.token
    router.push('/')
  } catch (error) {
    if (error instanceof ApiError) {
      if (error.code === 'AuthenticatingUser') {
        toast.add({
          severity: 'error',
          summary: 'Error',
          detail: 'Invalid credentials',
          life: 3000,
        })
        return
      }
    }
    console.error(error)
  }
}
const signup = async () => {
  try {
    await register(name.value, email.value, password.value)
    loginLocal(email.value, password.value)
  } catch (error) {
    if (error instanceof ApiError) {
      if (error.code === 'ObjectAlreadyExists') {
        console.error('User already exists')
        return
      }
    }
    console.error(error)
  }
}
</script>

<template>
  <div class="center">
    <Transition name="slide-up" mode="out-in">
      <LoginCard
        v-if="!signupMode"
        v-model:email="email"
        v-model:password="password"
        v-model:name="name"
        type="login"
        @on-submit="login"
        @on-switch="signupMode = !signupMode"
      />
      <LoginCard
        v-else-if="signupMode"
        v-model:email="email"
        v-model:password="password"
        v-model:name="name"
        type="signup"
        @on-submit="signup"
        @on-switch="signupMode = !signupMode"
      />
    </Transition>
  </div>
</template>

<style scoped>
.center {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100vh;
}

.slide-up-enter-active,
.slide-up-leave-active {
  transition: all 0.25s ease-out;
}

.slide-up-enter-from {
  opacity: 0;
  transform: translateX(-30px);
}

.slide-up-leave-to {
  opacity: 0;
  transform: translateX(30px);
}
</style>
