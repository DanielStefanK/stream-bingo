<script lang="ts" setup>
import Card from 'primevue/card'
import Button from 'primevue/button'
import Divider from 'primevue/divider'
import FloatLabel from 'primevue/floatlabel'
import InputText from 'primevue/inputtext'

const login = async () => {}
</script>

<template>
  <div class="center">
    <Card style="width: 25rem; overflow: hidden">
      <template #header>
        <img
          alt="user header"
          src="https://primefaces.org/cdn/primevue/images/usercard.png"
        />
      </template>
      <template #title>Login</template>
      <template #content>
        <div class="cta-oauth">Login with you favorit provider:</div>
        <div class="oauth">
          <Button class="oauth-login-btn twitch" unstyled>
            <i class="pi pi-twitch"></i>
          </Button>
          <Button class="oauth-login-btn google" unstyled>
            <i class="pi pi-google"></i>
          </Button>
          <Button class="oauth-login-btn github" unstyled>
            <i class="pi pi-github"></i>
          </Button>
        </div>

        <Divider align="center" type="dotted">
          <b>Or</b>
        </Divider>

        <div class="inputs">
          <FloatLabel variant="in">
            <InputText
              id="email"
              v-model="value2"
              variant="filled"
              class="w-full"
            />
            <label for="email">e-mail</label>
          </FloatLabel>

          <FloatLabel variant="in">
            <InputText
              id="password"
              v-model="value2"
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
          <Button label="Sign up" variant="link" />
          <Button @click="login">
            <span>Login</span>
            <i class="pi pi-sign-in" style="margin-left: 0.5rem"></i>
          </Button>
        </div>
      </template>
    </Card>
  </div>
</template>

<style scoped>
.center {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100vh;
}
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
