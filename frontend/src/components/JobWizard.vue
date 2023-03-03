<script setup lang="ts">
import SearchableDropdown from '@/components/SearchableDropdown.vue'
import { computed, watch, reactive } from 'vue'
import { createProject } from '@/api/project'
import { useResourcesStore } from '@/stores/resources'
import { storeToRefs } from 'pinia'
import { useWizardHashDetect } from '@/composables/useWizardHashDetect'

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

const steps = [
  { name: 'Name Project' },
  { name: 'Add Hashlist' },
  { name: 'Configure Attack Settings' },
  { name: 'Review & Start Attack' }
].slice(props.firstStep ?? 0)

const attackModes = [
  { name: 'Dictionary' },
  { name: 'Dictionary + Mask' },
  { name: 'Mask + Dictionary' }
]

/*
 * User Inputs
 */
const inputs = reactive({
  projectName: '',
  projectDesc: '',

  hashlistName: '',
  hashType: '0',
  hashes: '',

  activeStep: 0
})

const hashesArr = computed(() => {
  return inputs.hashes.split(/\s+/).filter((x) => !!x)
})

const {
  detectButtonClass,
  detectButtonClick,
  detectButtonText,
  detectStatusText,
  suggestedHashTypes,
  isLoadingSuggestions
} = useWizardHashDetect(hashesArr)

watch(suggestedHashTypes, (newHashTypes) => {
  const types = newHashTypes?.possible_types
  if (!types || types.length == 0) {
    return
  }
  inputs.hashType = types.sort()[0].toString()
})

const filteredHashTypes = computed(() => {
  const suggested = suggestedHashTypes.value?.possible_types
  if (suggested != null) {
    return allHashTypes.value.filter((hashType) => suggested.includes(hashType.id))
  }

  return allHashTypes.value
})

const hashTypeOptionsToShow = computed(() =>
  filteredHashTypes.value.map((type) => ({
    value: type.id.toString(),
    text: `${type.id} - ${type.name}`
  }))
)

const selectedHashType = computed(() =>
  allHashTypes.value.find((x) => x.id.toString() === inputs.hashType)
)

async function saveUptoProject() {
  createProject(inputs.projectName, inputs.projectDesc)
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
        :class="index <= inputs.activeStep ? 'step-primary step' : 'step'"
      >
        {{ step.name }}
      </li>
    </ul>
    <div
      class="card min-w-max self-center bg-base-100 shadow-xl"
      style="min-width: 800px; max-width: 80%"
    >
      <div class="card-body">
        <h2 class="card-title mb-8 w-96 justify-center self-center text-center">
          Step {{ inputs.activeStep + 1 }}. {{ steps[inputs.activeStep].name }}
        </h2>

        <template v-if="inputs.activeStep == 0">
          <input
            v-model="inputs.projectName"
            type="text"
            placeholder="Project Name"
            class="input-bordered input w-full max-w-xs"
          />
          <input
            v-model="inputs.projectDesc"
            type="text"
            placeholder="Project Description (optional)"
            class="input-bordered input w-full max-w-xs"
          />

          <div class="mt-8 flex justify-between">
            <div class="flex justify-start">
              <button class="link" @click="saveUptoProject">Create empty project and finish</button>
            </div>

            <div class="card-actions justify-end">
              <button class="btn-primary btn" @click="inputs.activeStep++">Next</button>
            </div>
          </div>
        </template>

        <template v-if="inputs.activeStep == 1">
          <div class="form-control">
            <label class="label font-bold">
              <span class="label-text">Hashlist Name</span>
            </label>
            <input
              type="text"
              placeholder="Dumped Admin NTLM Hashes"
              v-model="inputs.hashlistName"
              class="input-bordered input w-full max-w-xs"
            />
            <hr class="my-8" />
            <label class="label font-bold">
              <span class="label-text">Hashes (one per line)</span>
            </label>
            <textarea
              placeholder="Hashes"
              class="hashes-input textarea-bordered textarea w-full"
              rows="12"
              v-model="inputs.hashes"
            ></textarea>
            <label class="label mt-4 font-bold">
              <span class="label-text">Hash Type ({{ filteredHashTypes.length }} options)</span>
            </label>

            <SearchableDropdown v-model="inputs.hashType" :options="hashTypeOptionsToShow" placeholder-text="Search for a hashtype..." />

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

            <div class="mt-8 flex justify-between">
              <div class="flex justify-start">
                <button class="link" @click="saveUptoHashlist">Save hashlist and finish</button>
              </div>
              <div class="card-actions justify-end">
                <button class="btn-ghost btn" @click="inputs.activeStep--">Previous</button>
                <button class="btn-primary btn" @click="inputs.activeStep++">Next</button>
              </div>
            </div>
          </div>
        </template>

        <template v-if="inputs.activeStep == 2">
          <div class="btn-group self-center">
            <input
              type="radio"
              name="options"
              :data-title="attackMode.name"
              class="btn"
              :key="attackMode.name"
              v-for="attackMode in attackModes"
            />
          </div>

          <div class="mt-8 flex justify-between">
            <div class="flex justify-start">
              <button class="link" @click="saveUptoAttack">Save attack and finish</button>
            </div>

            <div class="card-actions justify-end">
              <button class="btn btn-ghost" @click="inputs.activeStep--">Previous</button>
              <button class="btn btn-primary" @click="inputs.activeStep++">Next</button>
            </div>
          </div>
        </template>

        <template v-if="inputs.activeStep == 3">
          <table class="first-col-bold table w-full">
            <tbody>
              <tr>
                <td>Project Name</td>
                <td>{{ inputs.projectName }}</td>
              </tr>
              <tr>
                <td>Project Description</td>
                <td>{{ inputs.projectDesc }}</td>
              </tr>
              <tr>
                <td>Hashlist Name</td>
                <td>TODO</td>
              </tr>
              <tr>
                <td>Hashlist Type</td>
                <td>{{ selectedHashType?.id }} - {{ selectedHashType?.name }}</td>
              </tr>
              <tr>
                <td>Number of Hashes</td>
                <td>{{ hashesArr.length }}</td>
              </tr>
            </tbody>
          </table>

          <div class="mt-8 flex justify-between">
            <div class="flex justify-start">
              <button class="link" @click="saveUptoAttack">Save attack and finish</button>
            </div>

            <div class="card-actions justify-end">
              <button class="btn btn-ghost" @click="inputs.activeStep--">Previous</button>
              <button class="btn btn-success" @click="inputs.activeStep++">Start Attack</button>
            </div>
          </div>
        </template>
      </div>
    </div>
  </div>
</template>

<style scoped>
textarea.hashes-input {
  background-image: linear-gradient(to bottom, rgba(87, 87, 87, 0.05) 50%, transparent 50%);
  background-repeat: repeat-y;

  background-size: 100% 56px;

  line-height: 28px;
  font-size: 18px;

  padding: 0;
  padding-left: 3px;

  white-space: pre;
  resize: none;

  background-attachment: local;
}

table.first-col-bold tr > td:first-of-type {
  font-weight: bold;
}
</style>
