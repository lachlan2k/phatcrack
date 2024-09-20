<script setup lang="ts">
import { storeToRefs } from 'pinia'
import { computed, ref } from 'vue'
import { useToast } from 'vue-toastification'

import Modal from '@/components/Modal.vue'
import IconButton from '@/components/IconButton.vue'
import AttackSettings from '@/components/Wizard/AttackSettings.vue'
import EmptyTable from '@/components/EmptyTable.vue'
import ConfirmModal from '@/components/ConfirmModal.vue'
import PageLoading from '@/components/PageLoading.vue'

import { AttackTemplateSetType, AttackTemplateType } from '@/api/attackTemplate'

import { useAttackSettings } from '@/composables/useAttackSettings'
import { useToastError } from '@/composables/useToastError'

import { useAttackTemplatesStore } from '@/stores/attackTemplates'

import { Icons } from '@/util/icons'

const attackTemplatesStore = useAttackTemplatesStore()
const { templates, isFirstLoading } = storeToRefs(attackTemplatesStore)
attackTemplatesStore.load(true)

const isCreateModalOpen = ref(false)
const isEditModalOpen = ref(false)

const attackTemplateToEditId = ref('')

const newAttackTemplateName = ref('')
const editAttackTemplateName = ref('')

const { attackSettings: attackSettingsForCreate, asHashcatParams: createAsHashcatParams } = useAttackSettings()
const {
  attackSettings: attackSettingsForEdit,
  asHashcatParams: editAsHashcatParams,
  loadFromHashcatParams: loadEditFromHashcatParams
} = useAttackSettings()

async function onOpenEditAttackSettings(id: string) {
  const tmpl = attackTemplatesStore.byId(id)
  if (!tmpl) {
    toast.warning('Failed to open editor - template was null')
    return
  }

  editAttackTemplateName.value = tmpl.name
  attackTemplateToEditId.value = tmpl.id

  switch (tmpl.type) {
    case AttackTemplateType: {
      if (tmpl.hashcat_params == null) {
        toast.warning('Failed to open editor - settings were null')
        return
      }

      loadEditFromHashcatParams(tmpl.hashcat_params)
      isEditModalOpen.value = true
      break
    }

    case AttackTemplateSetType: {
      break
    }
  }
}

const createAttackTemplateValidationError = computed(() => {
  if (newAttackTemplateName.value.length < 3) {
    return 'Name must be 3 or more characters'
  }

  return null
})

const editAttackTemplateValidationError = computed(() => {
  if (editAttackTemplateName.value.length < 3) {
    return 'Name must be 3 or more characters'
  }

  return null
})

const isLoadingForm = ref(false)

const toast = useToast()
const { catcher } = useToastError()

async function onCreateAttackTemplate() {
  isLoadingForm.value = true
  try {
    // todo loading spin
    const created = await attackTemplatesStore.create({
      name: newAttackTemplateName.value,
      hashcat_params: createAsHashcatParams(0)
    })
    toast.success(`Created new attack tempalte ${created.name}`)
  } catch (e) {
    catcher(e, 'Failed to create attack template')
  } finally {
    isLoadingForm.value = false
  }
}

async function onDeleteAttackTemplate(id: string) {
  try {
    await attackTemplatesStore.delete(id)
  } catch (e) {
    catcher(e, 'Failed to delete attack template')
  }
}

async function onSaveAttackTemplate() {
  try {
    isLoadingForm.value = true
    await attackTemplatesStore.update(attackTemplateToEditId.value, {
      name: editAttackTemplateName.value,
      type: AttackTemplateType,
      hashcat_params: editAsHashcatParams(0)
    })
    isEditModalOpen.value = false
    toast.success('Saved!')
  } catch (e) {
    catcher(e, 'Failed to save attack template')
  } finally {
    isLoadingForm.value = false
  }
}
</script>

<template>
  <Modal v-model:isOpen="isCreateModalOpen">
    <h3 class="text-lg font-bold mr-12 mb-4">Create Attack Template</h3>

    <div class="form-control">
      <label class="label font-bold"><span class="label-text">Name</span></label>
      <input type="text" class="input input-bordered" v-model="newAttackTemplateName" placeholder="Big Wordlist Attack" />
    </div>

    <hr class="my-4" />

    <AttackSettings v-model="attackSettingsForCreate" />

    <hr class="my-4" />
    <div class="tooltip tooltip-left float-right" :data-tip="createAttackTemplateValidationError">
      <button
        class="btn btn-primary"
        :disabled="createAttackTemplateValidationError != null || isLoadingForm"
        @click="() => onCreateAttackTemplate()"
      >
        <span class="loading loading-spinner loading-md" v-if="isLoadingForm"></span>
        Create
      </button>
    </div>
  </Modal>

  <Modal v-model:isOpen="isEditModalOpen">
    <h3 class="text-lg font-bold mr-12 mb-4">Edit Attack Template</h3>

    <div class="form-control">
      <label class="label font-bold"><span class="label-text">Name</span></label>
      <input type="text" class="input input-bordered" v-model="editAttackTemplateName" placeholder="Big Wordlist Attack" />
    </div>

    <hr class="my-4" />

    <AttackSettings v-model="attackSettingsForEdit" />

    <hr class="my-4" />
    <div class="tooltip tooltip-left float-right" :data-tip="editAttackTemplateValidationError">
      <button
        class="btn btn-primary"
        :disabled="editAttackTemplateValidationError != null || isLoadingForm"
        @click="() => onSaveAttackTemplate()"
      >
        <span class="loading loading-spinner loading-md" v-if="isLoadingForm"></span>
        Save
      </button>
    </div>
  </Modal>

  <main class="w-full h-full p-4">
    <PageLoading v-if="isFirstLoading" />
    <div v-else>
      <h1 class="text-4xl font-bold">Attack Templates</h1>
      <div class="mt-6 flex flex-wrap gap-6">
        <div class="card min-w-[800px] bg-base-100 shadow-xl">
          <div class="card-body">
            <div class="flex flex-row justify-between">
              <h2 class="card-title">Attack Templates</h2>
              <div>
                <button class="btn btn-primary btn-sm ml-1" @click="() => (isCreateModalOpen = true)">New Template</button>
                <button class="btn btn-primary btn-sm ml-1" @click="() => (isCreateModalOpen = true)" :disabled="templates.length === 0">
                  New Template Set
                </button>
              </div>
            </div>
            <table class="table w-full">
              <thead>
                <tr>
                  <th>Type</th>
                  <th>Name</th>
                  <th>Actions</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="tmpl in templates" :key="tmpl.id">
                  <td>
                    <font-awesome-icon :icon="Icons.AttackTemplate" v-if="tmpl.type === AttackTemplateType" />
                    <font-awesome-icon :icon="Icons.AttackTemplateSet" v-if="tmpl.type === AttackTemplateSetType" />
                  </td>
                  <td>
                    {{ tmpl.name }}
                  </td>
                  <td>
                    <IconButton @click="() => onOpenEditAttackSettings(tmpl.id)" :icon="Icons.Edit" color="primary" tooltip="Edit" />
                    <ConfirmModal @on-confirm="() => onDeleteAttackTemplate(tmpl.id)">
                      <IconButton :icon="Icons.Delete" color="error" tooltip="Delete" />
                    </ConfirmModal>
                  </td>
                </tr>
              </tbody>
            </table>
            <EmptyTable v-if="templates.length == 0" text="No Attack Templates Yet" :icon="Icons.AttackTemplate" />
          </div>
        </div>
      </div>
    </div>
  </main>
</template>
