<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { createProject } from '@/api/project'
import { detectHashType } from '@/api/resources'
import { useResourcesStore } from '@/stores/resources'
import { storeToRefs } from 'pinia'
import { useApi } from '@/composables/useApi'

/*
 * Props
 */

const props = defineProps<{
  // Set to 0 for full wizard, 1 if project is already made, 2 if hashlist is already made...
  firstStep?: number
  existingProjectID?: string
  existingHashlistId?: string
}>()

const resourcesStore = useResourcesStore()
const { hashTypes: allHashTypes } = storeToRefs(resourcesStore)
resourcesStore.loadHashTypes()

/*
 * User Inputs
 */

// Project parameters
const projectName = ref('')
const projectDesc = ref('')

// Hashlist parameters
const hashTypeFilter = ref('')
const hashType = ref(0)
const hashes = ref('')

const hashesArr = computed(() => {
  return hashes.value.split(/\s+/).filter((x) => !!x)
})

/*
 * Computed/helpers/state
 */

// Wizard steps
const activeStep = ref(0)

const steps = [
  { name: 'Name Project' },
  { name: 'Add Hashlist' },
  { name: 'Configure Attack Settings' },
  { name: 'Review & Start Attack' }
].slice(props.firstStep ?? 0)

// Hash type list
const {
  fetchData: fetchHashTypeSuggestions,
  isLoading: isLoadingSuggestions,
  data: suggestedHashTypes
} = useApi(() => detectHashType(hashesArr.value[0]), { immediate: false })

function detectButtonClick() {
  // Reset
  console.log('doing the thing')
  if (suggestedHashTypes.value != null) {
    suggestedHashTypes.value = null
    return
  }

  fetchHashTypeSuggestions()
}

const detectStatusText = computed(() => {
  if (suggestedHashTypes.value == null) {
    return ''
  }

  if (suggestedHashTypes.value.possible_types.length > 0) {
    return `Filtered down to ${suggestedHashTypes.value.possible_types.length} possible hash types`
  }

  return 'No suggestions found, check your hashes are valid'
})

const detectButtonClass = computed(() => {
  if (isLoadingSuggestions.value) {
    return 'btn-secondary'
  }

  if (suggestedHashTypes.value != null) {
    return ''
  }

  return 'btn-primary'
})

const detectButtonText = computed(() => {
  if (isLoadingSuggestions.value) {
    return 'Loading suggestions...'
  }

  if (suggestedHashTypes.value != null) {
    return 'Reset Filter'
  }

  return 'Detect hash type'
})

const filteredHashTypes = computed(() => {
  const suggested = suggestedHashTypes.value?.possible_types
  if (suggested != null) {
    return allHashTypes.value.filter((hashType) => suggested.includes(hashType.id))
  }

  const filterStr = hashTypeFilter.value.toLowerCase()
  if (hashTypeFilter.value === '') {
    return allHashTypes.value
  }

  return allHashTypes.value.filter(
    (hashType) =>
      hashType.id.toString().includes(filterStr) || hashType.name.toLowerCase().includes(filterStr)
  )
})

watch(suggestedHashTypes, (newHashTypes) => {
  const types = newHashTypes?.possible_types
  if (!types || types.length == 0) {
    return
  }
  hashType.value = types.sort()[0]
})

async function saveUptoProject() {
  createProject(projectName.value, projectDesc.value)
}

async function saveUptoHashlist() {
  await saveUptoProject()
}

async function saveUptoAttack() {
  await saveUptoHashlist()
}
</script>

<template>
  <div class="mt-6 flex flex-col flex-wrap gap-6">
    <ul class="steps my-8">
      <li
        v-for="(step, index) in steps"
        :key="index"
        :class="index <= activeStep ? 'step-primary step' : 'step'"
      >
        {{ step.name }}
      </li>
    </ul>
    <div
      class="card min-w-max self-center bg-base-100 shadow-xl"
      style="min-width: 800px; max-width: 80%"
    >
      <div class="card-body">
        <h2 class="card-title mb-8 w-96">
          Step {{ activeStep + 1 }}. {{ steps[activeStep].name }}
        </h2>

        <template v-if="activeStep == 0">
          <input
            v-model="projectName"
            type="text"
            placeholder="Project Name"
            class="input-bordered input w-full max-w-xs"
          />
          <input
            v-model="projectDesc"
            type="text"
            placeholder="Project Description (optional)"
            class="input-bordered input w-full max-w-xs"
          />

          <div class="mt-8 flex justify-between">
            <div class="flex justify-start">
              <button class="link" @click="saveUptoProject">Create empty project and finish</button>
            </div>

            <div class="card-actions justify-end">
              <button class="btn-primary btn" @click="activeStep++">Next</button>
            </div>
          </div>
        </template>

        <template v-if="activeStep == 1">
          <div class="form-control">
            <label class="label">
              <span class="label-text">Filter Hash Type</span>
            </label>
            <input
              type="text"
              placeholder="Search hash types"
              v-model="hashTypeFilter"
              class="input-bordered input w-full max-w-xs"
            />
            <label class="label mt-4">
              <span class="label-text">Hash Types ({{ filteredHashTypes.length }})</span>
            </label>
            <div>
              <select class="input-bordered input w-full max-w-xs" v-model="hashType">
                <option v-for="thisHashType in filteredHashTypes" :key="thisHashType.id" :value="thisHashType.id">
                  {{ thisHashType.id }} - {{ thisHashType.name }}
                </option>
              </select>
            </div>

            <div class="my-4">
              <button
                class="btn-sm btn"
                :class="detectButtonClass"
                :disabled="isLoadingSuggestions || hashesArr.length == 0"
                @click="detectButtonClick"
              >
                {{ detectButtonText }}
              </button>
              <span class="ml-2">{{ detectStatusText }}</span>
            </div>

            <label class="label">
              <span class="label-text">Hashes (one per line)</span>
            </label>
            <textarea
              placeholder="Hashes"
              class="textarea-bordered textarea max-w-xs hashes-input"
              v-model="hashes"
            ></textarea>

            <div class="mt-8 flex justify-between">
              <div class="flex justify-start">
                <button class="link" @click="saveUptoHashlist">Save hashlist and finish</button>
              </div>
              <div class="card-actions justify-end">
                <button class="btn-ghost btn" @click="activeStep--">Previous</button>
                <button class="btn-primary btn" @click="activeStep++">Next</button>
              </div>
            </div>
          </div>
        </template>

        <template v-if="activeStep == 2">
          <input type="text" placeholder="Hash Type" class="input-bordered input w-full max-w-xs" />
          <textarea placeholder="Hashes" class="textarea-bordered textarea max-w-xs"></textarea>

          <div class="mt-8 flex justify-between">
            <div class="flex justify-start">
              <button class="link" @click="saveUptoAttack">Save attack and finish</button>
            </div>

            <div class="card-actions justify-end">
              <button class="btn-ghost btn" @click="activeStep--">Previous</button>
              <button class="btn-primary btn" @click="activeStep++">Next</button>
            </div>
          </div>
        </template>

        <template v-if="activeStep == 3">
          <p><strong>Project Name</strong> {{ projectName }}</p>

          <div class="mt-8 flex justify-between">
            <div class="flex justify-start">
              <button class="link" @click="saveUptoAttack">Save attack and finish</button>
            </div>

            <div class="card-actions justify-end">
              <button class="btn-ghost btn" @click="activeStep--">Previous</button>
              <button class="btn-success btn" @click="activeStep++">Start Attack</button>
            </div>
          </div>
        </template>
      </div>
    </div>
  </div>
</template>

<style scoped>
textarea.hashes-input {
  background-image: linear-gradient(to bottom, rgba(0,0,0,.05) 50%, transparent 50%);
  background-size: 100% 56px;
  line-height: 28px;
  /* line-height: 3em; */
  background-repeat: repeat-y;
  /* background-position-y: -1em; */
  padding: 0;
  padding-left: 3px;
  font-size: 18px;
}
</style>