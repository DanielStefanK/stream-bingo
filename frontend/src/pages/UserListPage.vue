<script setup lang="ts">
import { ref, watch } from 'vue'
import MainLayout from '../layout/MainLayout.vue'
import {
  deleteUser,
  loadUsers,
  toggleUserActive,
  type UserWithAllFields,
} from '../services/admin'
import { Card } from 'primevue'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import TableAction from '../components/TableAction.vue'

const page = ref(1)
const perPage = ref(10)
const total = ref(0)
const users = ref<UserWithAllFields[]>([])
const loading = ref(false)
const term = ref('')
const loadingAction = ref<Record<string | number, boolean>>({})

const fetchPage = async () => {
  loading.value = true
  try {
    const resp = await loadUsers(page.value, perPage.value, term.value)
    users.value = resp.data
    total.value = resp.total
  } catch (error) {
    console.error(error)
  } finally {
    loading.value = false
  }
}

const toggleActive = async (id: number, active: boolean) => {
  loadingAction.value[id] = true
  try {
    await toggleUserActive(id, active)
    await fetchPage()
  } catch (error) {
    console.error(error)
  } finally {
    loadingAction.value[id] = false
  }
}

const deleteUserWithId = async (id: number) => {
  loadingAction.value[id] = true
  try {
    await deleteUser(id)
    await fetchPage()
  } catch (error) {
    console.error(error)
  } finally {
    loadingAction.value[id] = false
  }
}

watch([page, perPage, term], fetchPage, { immediate: true })
</script>

<template>
  <MainLayout>
    <template #default>
      <Card>
        <template #content>
          <DataTable :value="users" paginator :rows="10" :total-records="total">
            <Column field="id" header="ID"></Column>
            <Column field="avatar" header="Avatar">
              <template #body="slotProps">
                <img
                  v-if="slotProps.data.avatar"
                  :alt="slotProps.data.name"
                  :src="slotProps.data.avatar"
                  width="32"
                  style="vertical-align: middle"
                />
              </template>
            </Column>
            <Column field="name" header="Name" style="width: 20%"></Column>
            <Column field="email" header="Email" style="width: 20%"></Column>
            <Column
              field="provider"
              header="Provider"
              style="width: 20%"
            ></Column>
            <Column field="admin" header="Admin" style="width: 10%"></Column>
            <Column field="active" header="Active" style="width: 10%"></Column>

            <Column header="Actions" style="width: 10%">
              <template #body="{ data: { id, active } }">
                <TableAction
                  :loading="loadingAction[id]"
                  :items="[
                    {
                      label: active ? 'Deactivate' : 'Activate',
                      icon: `pi pi-${active ? 'ban' : 'check-circle'}`,
                      command: () => {
                        toggleActive(id, !active)
                      },
                    },
                    {
                      label: 'Delete',
                      icon: 'pi pi-trash',
                      command: () => {
                        deleteUserWithId(id)
                      },
                    },
                  ]"
                />
              </template>
            </Column>
          </DataTable>
        </template>
      </Card>
    </template>
  </MainLayout>
</template>
