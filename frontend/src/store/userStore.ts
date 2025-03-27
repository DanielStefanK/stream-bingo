import { defineStore } from 'pinia'
import { fetchUser } from '../services/auth'

export const useUserStore = defineStore('userSTore', {
  state: () => ({
    user: null as User | null,
    token: '',
    state: 'uninitialized' as 'uninitialized' | 'loading' | 'loaded',
    userPromise: null as Promise<User> | null,
  }),
  getters: {
    userLoad(state) {
      if (state.userPromise) return state.userPromise
      if (state.user) return state.user
      if (!localStorage.getItem('token')) return null
      if (state.state === 'uninitialized') {
        state.state = 'loading'
        state.userPromise = fetchUser()
        state.userPromise
          .then((user) => {
            state.user = user
            state.state = 'loaded'
          })
          .catch(() => {
            state.state = 'loaded'
          })
        return state.userPromise
      }
      return null
    },
  },
})

export interface User {
  id: number
  name: string
  email: string
  avatarURL?: string
  provider: string
  admin: boolean
}
