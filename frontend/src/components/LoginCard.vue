<script setup lang="ts">
import Card from 'primevue/card'
import Button from 'primevue/button'
import Divider from 'primevue/divider'
import FloatLabel from 'primevue/floatlabel'
import InputText from 'primevue/inputtext'

import { getLoginUrl, providers } from '../services/auth'

const e = defineEmits(['onSubmit', 'onSwitch'])
const email = defineModel<string>('email')
const password = defineModel<string>('password')
const name = defineModel<string>('name')
defineProps<{ type: 'login' | 'signup' }>()
</script>

<template>
  <Card style="width: 25rem; overflow: hidden">
    <template #header>
      <img
        alt="user header"
        src="https://primefaces.org/cdn/primevue/images/usercard.png"
      />
    </template>
    <template #title>{{ type === 'signup' ? 'Sign-up' : 'Login' }}</template>
    <template #content>
      <div class="cta-oauth">
        {{ type === 'signup' ? 'Sign-up' : 'Login' }} with you favorit provider:
      </div>
      <div class="oauth">
        <Button
          v-for="provider in providers"
          :key="provider.id"
          as="a"
          :class="`oauth-login-btn ${provider.id}`"
          :href="getLoginUrl(provider.id)"
          unstyled
        >
          <i :class="`pi pi-${provider.id}`"></i>
        </Button>
      </div>

      <Divider align="center" type="dotted">
        <b>Or</b>
      </Divider>

      <div class="inputs">
        <FloatLabel v-if="type === 'signup'" variant="in">
          <InputText id="name" v-model="name" variant="filled" class="w-full" />
          <label for="name">name</label>
        </FloatLabel>
        <FloatLabel variant="in">
          <InputText
            id="email"
            v-model="email"
            variant="filled"
            class="w-full"
          />
          <label for="email">e-mail</label>
        </FloatLabel>

        <FloatLabel variant="in">
          <InputText
            id="password"
            v-model="password"
            variant="filled"
            class="w-full"
            type="password"
          />
          <label for="password">password</label>
        </FloatLabel>
      </div>
    </template>
    <template #footer>
      <div class="actions">
        <Button
          :label="type === 'login' ? 'Sign-up' : 'Login'"
          variant="link"
          @click="e('onSwitch')"
        />
        <Button @click="e('onSubmit')">
          <span>{{ type === 'signup' ? 'Sign-up' : 'Login' }}</span>
          <i class="pi pi-sign-in" style="margin-left: 0.5rem"></i>
        </Button>
      </div>
    </template>
  </Card>
</template>

<style scoped>
.actions {
  display: flex;
  justify-content: end;
  gap: 1rem;
  margin-top: 1rem;
}
.inputs {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.oauth {
  display: flex;
  justify-content: space-evenly;
  gap: 1rem;
}

.cta-oauth {
  margin-bottom: 1rem;
  color: var(--p-card-subtitle-color);
}

.oauth-login-btn {
  display: flex;
  align-items: center;
  gap: 8px;
  color: white;
  font-weight: bold;
  padding: 10px 16px;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
  transition:
    background 0.2s ease-in-out,
    transform 0.1s ease-in-out;
  text-decoration: none;
}

.oauth-login-btn:hover {
  background-color: #772ce8;
  transform: scale(1.05);
}

.oauth-login-btn:active {
  transform: scale(0.98);
}

.twitch {
  background-color: #6441a4;
}
.twitch:hover {
  background-color: #4c2d7e;
}

.google {
  background-color: #db4437;
}
.google:hover {
  background-color: #c1351d;
}
.github {
  background-color: #333;
}
.github:hover {
  background-color: #222;
}
</style>
