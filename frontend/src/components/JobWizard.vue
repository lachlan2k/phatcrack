<script setup lang="ts">
import FileSelectModal from '@/components/Wizard/FileSelectModal.vue'
import SearchableDropdown from '@/components/SearchableDropdown.vue'
import HrOr from '@/components/HrOr.vue'
import { computed, watch, reactive } from 'vue'
import { createProject } from '@/api/project'
import { useResourcesStore } from '@/stores/resources'
import { storeToRefs } from 'pinia'
import { useWizardHashDetect } from '@/composables/useWizardHashDetect'
import { useProjectsStore } from '@/stores/projects'
import { useApi } from '@/composables/useApi'
import { getAllRulefiles, getAllWordlists } from '@/api/lists'

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

const projectsStore = useProjectsStore()
const { projects } = storeToRefs(projectsStore)
projectsStore.load()

const projectSelectOptions = computed(() => [
  { value: '', text: 'Create new project' },
  ...projects.value.map((project) => ({
    value: project.id,
    text: project.name
  }))
])

const { data: allWordlists } = useApi(getAllWordlists)
const { data: allRulefiles } = useApi(getAllRulefiles)

const steps = [
  { name: 'Choose or Create Project' },
  { name: 'Add Hashlist' },
  { name: 'Configure Attack Settings' },
  { name: 'Review & Start Attack' }
].slice(props.firstStep ?? 0)

const attackModes = [
  { name: 'Wordlist', value: 0 },
  { name: 'Combinator', value: 1 },
  { name: 'Brute-force/Mask', value: 3 },
  { name: 'Wordlist + Mask', value: 6 },
  { name: 'Mask + Wordlist', value: 7 }
]

/*
 * User Inputs
 */
const inputs = reactive({
  projectName: '',
  projectDesc: '',
  selectedProjectId: '',

  hashlistName: '',
  hashType: '0',
  hashes: '',

  attackMode: 0,
  selectedWordlists: [] as string[],
  selectedRulefiles: [] as string[],
  optimizedKernels: false,
  slowCandidates: false,
  enableLoopback: true,

  activeStep: 2
})

function toggledSelectWordlist(id: string) {
  if (inputs.selectedWordlists.includes(id)) {
    inputs.selectedWordlists = inputs.selectedWordlists.filter(x => x != id)
  } else {
    inputs.selectedWordlists.push(id)
  }
}

function toggledSelectRulefile(id: string) {
  if (inputs.selectedRulefiles.includes(id)) {
    inputs.selectedRulefiles = inputs.selectedRulefiles.filter(x => x != id)
  } else {
    inputs.selectedRulefiles.push(id)
  }
}

watch(
  () => inputs.projectName,
  (newProjName) => {
    if (newProjName != '') {
      inputs.selectedProjectId = ''
    }
  }
)

watch(
  () => inputs.selectedProjectId,
  (newSelectedProj) => {
    if (newSelectedProj != '') {
      inputs.projectName = ''
      inputs.projectDesc = ''
    }
  }
)

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

/*
 * Step validations
 */
const projectStepValidationError = computed(() => {
  if (inputs.projectName == '' && inputs.selectedProjectId == '') {
    return 'Please select an existing project or input a new project name'
  }
  return null
})

const hashlistStepValidationError = computed(() => {
  if (inputs.hashlistName == '') {
    return 'Please name the hashlist'
  }

  if (hashesArr.value.length == 0) {
    return 'Please input at least one hash'
  }
  return null
})

/*
 * API Helpers
 */
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
      <li v-for="(step, index) in steps" :key="index" :class="index <= inputs.activeStep ? 'step-primary step' : 'step'">
        {{ step.name }}
      </li>
    </ul>
    <div class="card min-w-max self-center bg-base-100 shadow-xl" style="min-width: 800px; max-width: 80%">
      <div class="card-body">
        <h2 class="card-title mb-8 w-96 justify-center self-center text-center">
          Step {{ inputs.activeStep + 1 }}. {{ steps[inputs.activeStep].name }}
        </h2>

        <!-- Create/Select Project -->
        <template v-if="inputs.activeStep == 0">
          <div class="form-control">
            <label class="label font-bold">
              <span class="label-text">Choose Project</span>
            </label>
            <SearchableDropdown v-model="inputs.selectedProjectId" :options="projectSelectOptions"
              placeholderText="Select existing project..." class="max-w-xs" />

            <HrOr class="my-4" />

            <label class="label font-bold">
              <span class="label-text">New Project Name</span>
            </label>
            <input v-model="inputs.projectName" type="text" placeholder="12345 Example Corp"
              class="input-bordered input w-full max-w-xs" />
            <label class="label mt-4 font-bold">
              <span class="label-text">New Project Description (optional)</span>
            </label>
            <input v-model="inputs.projectDesc" type="text" placeholder="Internal Network Pentest"
              class="input-bordered input w-full max-w-xs" />
            <div class="mt-8 flex justify-between">
              <div class="flex justify-start">
                <button class="link" @click="saveUptoProject" v-if="inputs.projectName != ''">
                  Create empty project and finish
                </button>
              </div>
              <div class="card-actions justify-end">
                <div class="tooltip" :data-tip="projectStepValidationError">
                  <button class="btn-primary btn" @click="inputs.activeStep++"
                    :disabled="projectStepValidationError != null">
                    Next
                  </button>
                </div>
              </div>
            </div>
          </div>
        </template>

        <!-- Create Hashlist -->
        <template v-if="inputs.activeStep == 1">
          <div class="form-control">
            <label class="label font-bold">
              <span class="label-text">Hashlist Name</span>
            </label>
            <input type="text" placeholder="Dumped Admin NTLM Hashes" v-model="inputs.hashlistName"
              class="input-bordered input w-full max-w-xs" />
            <hr class="my-8" />
            <label class="label font-bold">
              <span class="label-text">Hashes (one per line)</span>
            </label>
            <textarea placeholder="Hashes"
              class="hashes-input textarea-bordered textarea w-full font-mono focus:outline-none" rows="12"
              v-model="inputs.hashes"></textarea>
            <label class="label mt-4 font-bold">
              <span class="label-text">Hash Type ({{ filteredHashTypes.length }} options)</span>
            </label>

            <SearchableDropdown v-model="inputs.hashType" :options="hashTypeOptionsToShow"
              placeholder-text="Search for a hashtype..." />

            <div class="my-4">
              <button class="btn-sm btn" :class="detectButtonClass"
                :disabled="isLoadingSuggestions || hashesArr.length == 0" @click="detectButtonClick">
                {{ detectButtonText }}
              </button>
              <span class="ml-2">{{ detectStatusText }}</span>
            </div>

            <div class="mt-8 flex justify-between">
              <div class="flex justify-start">
                <button class="link" @click="saveUptoHashlist" v-if="hashlistStepValidationError == null">
                  Save hashlist and finish
                </button>
              </div>
              <div class="card-actions justify-end">
                <button class="btn-ghost btn" @click="inputs.activeStep--">Previous</button>
                <div class="tooltip" :data-tip="hashlistStepValidationError">
                  <button class="btn-primary btn" @click="inputs.activeStep++"
                    :disabled="hashlistStepValidationError != null">
                    Next
                  </button>
                </div>
              </div>
            </div>
          </div>
        </template>

        <!-- Attack settings -->
        <template v-if="inputs.activeStep == 2">
          <div class="btn-group self-center">
            <input type="radio" name="options" :data-title="attackMode.name" class="btn" :key="attackMode.value"
              :value="attackMode.value" v-model="inputs.attackMode" v-for="attackMode in attackModes" />
          </div>

          <!-- <FileSelectModal v-model="inputs.selectedRulefiles" open-button-text="Select Wordlists" :allow-multiple="true" :files="allRulefiles.rulefiles" v-if="allRulefiles != null" /> -->
          <!-- <FileSelectModal v-model="inputs.selectedRulefiles" open-button-text="Select Rulefiles" :allow-multiple="true" :files="allRulefiles.rulefiles" v-if="allRulefiles != null" /> -->

          <div class="my-2"></div>

          <label class="label font-bold">Select Wordlist(s)</label>
          <table class="table w-full">
            <tbody>
              <tr>
                <td>Select</td>
                <td>Name</td>
                <td>Number of lines</td>
              </tr>
              <tr v-for="file in allWordlists?.wordlists ?? []" :key="file.id" @click="toggledSelectWordlist(file.id)"
                class="cursor-pointer">
                <td>
                  <input type="checkbox" class="checkbox checkbox-primary checkbox-xs align-middle"
                    :checked="inputs.selectedWordlists.includes(file.id)">
                </td>
                <td>{{ file.name }}</td>
                <td>{{ file.lines }}</td>
              </tr>
            </tbody>
          </table>

          <hr class="my-4">

          <label class="label font-bold">Select Rule File(s)</label>

          <table class="table w-full">
            <tbody>
              <tr>
                <td>Select</td>
                <td>Name</td>
                <td>Number of lines</td>
              </tr>
              <tr v-for="file in allRulefiles?.rulefiles ?? []" :key="file.id" @click="toggledSelectRulefile(file.id)"
                class="cursor-pointer">
                <td>
                  <input type="checkbox" class="checkbox checkbox-primary checkbox-xs align-middle"
                    :checked="inputs.selectedRulefiles.includes(file.id)">
                </td>
                <td>{{ file.name }}</td>
                <td>{{ file.lines }}</td>
              </tr>
            </tbody>
          </table>

            <hr class="my-4">


          <label class="label font-bold">Additional Options</label>
          <div class="pl-3">
            <label class="label cursor-pointer justify-start">
                <input type="checkbox" v-model="inputs.enableLoopback" class="checkbox checkbox-xs checkbox-primary" />
                <span><span class="label-text font-bold ml-4">Loopback?</span> (--loopback)</span>
              </label>
            <label class="label cursor-pointer justify-start">
              <input type="checkbox" v-model="inputs.optimizedKernels" class="checkbox checkbox-xs checkbox-primary" />
              <span><span class="label-text font-bold ml-4">Optimized Kernels?</span> (-O)</span>
            </label>
            <label class="label cursor-pointer justify-start">
                <input type="checkbox" v-model="inputs.slowCandidates" class="checkbox checkbox-xs checkbox-primary" />
                <span><span class="label-text font-bold ml-4">Slow Candidates?</span> (-S)</span>
              </label>
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

        <!-- Review/start -->
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

  background-size: 100% 60px;

  line-height: 30px;
  font-size: 18px;

  padding: 0;
  padding-left: 3px;

  white-space: pre;
  resize: none;

  background-attachment: local;
}

table.first-col-bold tr>td:first-of-type {
  font-weight: bold;
}
</style>
