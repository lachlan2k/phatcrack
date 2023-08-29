<script setup lang="ts">
import { AxiosError } from 'axios'
import AttackConfigDetails from '@/components/AttackConfigDetails.vue'
import HashlistInputs from './HashlistInputs.vue'
import SearchableDropdown from '@/components/SearchableDropdown.vue'
import MaskInput from './MaskInput.vue'
import WordlistSelect from '@/components/Wizard/ListSelect.vue'
import HrOr from '@/components/HrOr.vue'
import { computed, watch, reactive } from 'vue'
import {
  createHashlist,
  createProject,
  createAttack,
  startAttack,
  getProject,
  getHashlist
} from '@/api/project'
import { storeToRefs } from 'pinia'
import { useProjectsStore } from '@/stores/projects'
import { useApi } from '@/composables/useApi'
import { getAllRulefiles, getAllWordlists } from '@/api/listfiles'
import type { AttackDTO, HashcatParams, HashlistCreateResponseDTO, ProjectDTO } from '@/api/types'
import { useToast } from 'vue-toastification'
import { attackModes } from '@/util/hashcat'
import { useResourcesStore } from '@/stores/resources'

/*
 * Props
 */
const props = withDefaults(
  defineProps<{
    // Set to 0 for full wizard, 1 if project is already made, 2 if hashlist is already made...
    firstStep?: number
    existingProjectId?: string
    existingHashlistId?: string
  }>(),
  {
    firstStep: 0
  }
)

interface StartEmitDetails {
  projectId: string
  hashlistId: string
  attackId: string
}

const emit = defineEmits<{
  (e: 'successfulStart', details: StartEmitDetails): void
}>()

const projectsStore = useProjectsStore()
const { projects } = storeToRefs(projectsStore)
projectsStore.load()

const resourcesStore = useResourcesStore()
const { hashTypes: allHashTypes } = storeToRefs(resourcesStore)
resourcesStore.loadHashTypes()

const projectSelectOptions = computed(() => [
  { value: '', text: 'Create new project 🖋' },
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
]

enum StepIndex {
  Project = 0,
  Hashlist = 1,
  Attack = 2,
  Review = 3
}

const stepsToDisplay = steps.slice(props.firstStep)

/*
 * User Inputs
 */
const inputs = reactive({
  projectName: '',
  projectDesc: '',
  selectedProjectId: props.existingProjectId ?? '',

  hashlistName: '',
  selectedHashlistId: props.existingHashlistId ?? '',
  hashType: '0',
  hashes: '',

  attackMode: 0,
  selectedWordlists: [] as string[],
  selectedRulefiles: [] as string[],
  mask: '',
  optimizedKernels: false,
  slowCandidates: false,
  enableLoopback: true,

  combinatorLeft: [] as string[],
  combinatorRight: [] as string[],

  isDistributed: true,

  activeStep: props.firstStep
})

// TODO: refactor so that selectedWordlists isn't the source of truth
watch(
  () => inputs.combinatorLeft,
  (newLeft) => (inputs.selectedWordlists = [...newLeft, ...inputs.combinatorRight])
)
watch(
  () => inputs.combinatorRight,
  (newRight) => (inputs.selectedWordlists = [...inputs.combinatorLeft, ...newRight])
)

// If a user starts typing in a new project name, then de-select existing project
watch(
  () => inputs.projectName,
  (newProjName) => {
    if (newProjName != '') {
      inputs.selectedProjectId = ''
    }
  }
)

// If a user selects an existing project, remove the project name they've typed
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

const toast = useToast()

/*
 * API Helpers
 */
async function saveOrGetProject(): Promise<ProjectDTO> {
  try {
    if (inputs.selectedProjectId) {
      const proj = await getProject(inputs.selectedProjectId)
      return proj
    }

    const proj = await createProject(inputs.projectName, inputs.projectDesc)

    inputs.selectedProjectId = proj.id

    toast.success(`Created project "${inputs.projectName}"!`)
    return proj
  } catch (err: any) {
    toast.warning('Failed to create project' + err.message)
    // Throw up so our caller knows an error happened
    throw err
  }
}

async function saveOrGetHashlist(): Promise<HashlistCreateResponseDTO> {
  const proj = await saveOrGetProject()

  try {
    if (inputs.selectedHashlistId) {
      const hashlist = await getHashlist(inputs.selectedHashlistId)
      return hashlist
    }

    const hashlist = await createHashlist({
      project_id: proj.id,
      name: inputs.hashlistName,
      hash_type: Number(inputs.hashType),
      input_hashes: hashesArr.value,
      has_usernames: false
    })

    console.log('Created hashlist', hashlist)
    toast.success(`Created hashlist "${inputs.hashlistName}"!`)

    inputs.selectedHashlistId = hashlist.id

    return hashlist
  } catch (err: any) {
    toast.warning('Failed to create hashlist' + err.message)
    throw err
  }
}

function makeHashcatParams(): HashcatParams {
  return {
    attack_mode: inputs.attackMode,
    hash_type: Number(inputs.hashType),

    // TODO mask inputs
    // Also TODO: how does combinator work?
    mask: inputs.mask,
    mask_increment: true,
    mask_increment_min: 0,
    mask_increment_max: 0,
    mask_custom_charsets: [],
    mask_sharded_charset: '',

    wordlist_filenames: inputs.selectedWordlists,
    rules_filenames: inputs.selectedRulefiles,

    optimized_kernels: inputs.optimizedKernels,
    slow_candidates: inputs.slowCandidates,

    additional_args: [],
    skip: 0,
    limit: 0
  }
}

const computedHashcatParams = computed(() => {
  return makeHashcatParams()
})

async function saveUptoAttack(): Promise<AttackDTO> {
  const hashlist = await saveOrGetHashlist()

  try {
    const attack = await createAttack({
      hashlist_id: hashlist.id,
      hashcat_params: makeHashcatParams(),
      is_distributed: inputs.isDistributed
    })
    toast.success('Created attack!')
    return attack
  } catch (err: any) {
    toast.warning('Failed to create attack' + err.message)
    throw err
  }
}

async function saveAndStartAttack() {
  const attack = await saveUptoAttack()
  try {
    await startAttack(attack.id)
    emit('successfulStart', {
      projectId: inputs.selectedProjectId,
      hashlistId: inputs.selectedHashlistId,
      attackId: attack.id
    })
    toast.success('Started attack!')
  } catch (e: any) {
    let errorString = 'Unknown Error'
    if (e instanceof AxiosError) {
      errorString = e.response?.data?.message
    } else if (e instanceof Error) {
      errorString = e.message
    }

    toast.warning('Failed to start: ' + errorString)
  }
}
</script>

<template>
  <div class="mt-6 flex flex-col flex-wrap gap-6">
    <ul class="steps my-1">
      <li
        v-for="(step, index) in stepsToDisplay"
        :key="index"
        :class="index + props.firstStep <= inputs.activeStep ? 'step step-primary' : 'step'"
      >
        {{ step.name }}
      </li>
    </ul>
    <div
      class="card min-w-max self-center bg-base-100 shadow-xl"
      style="min-width: 800px; max-width: 80%"
    >
      <div class="card-body">
        <h2 class="card-title mb-4 w-96 justify-center self-center text-center">
          Step {{ inputs.activeStep + 1 - props.firstStep }}. {{ steps[inputs.activeStep].name }}
        </h2>

        <!-- Create/Select Project -->
        <template v-if="inputs.activeStep == StepIndex.Project">
          <div class="form-control">
            <label class="label font-bold">
              <span class="label-text">Choose Project</span>
            </label>
            <SearchableDropdown
              v-model="inputs.selectedProjectId"
              :options="projectSelectOptions"
              placeholderText="Select existing project..."
              class="max-w-xs"
            />

            <HrOr class="my-4" />

            <label class="label font-bold">
              <span class="label-text">New Project Name</span>
            </label>
            <input
              v-model="inputs.projectName"
              type="text"
              placeholder="12345 Example Corp"
              class="input input-bordered w-full max-w-xs"
            />
            <label class="label mt-4 font-bold">
              <span class="label-text">New Project Description (optional)</span>
            </label>
            <input
              v-model="inputs.projectDesc"
              type="text"
              placeholder="Internal Network Pentest"
              class="input input-bordered w-full max-w-xs"
            />
            <div class="mt-8 flex justify-between">
              <div class="flex justify-start">
                <button class="link" @click="saveOrGetProject" v-if="inputs.projectName != ''">
                  Create empty project and finish
                </button>
              </div>
              <div class="card-actions justify-end">
                <div class="tooltip" :data-tip="projectStepValidationError">
                  <button
                    class="btn btn-primary"
                    @click="inputs.activeStep++"
                    :disabled="projectStepValidationError != null"
                  >
                    Next
                  </button>
                </div>
              </div>
            </div>
          </div>
        </template>

        <!-- Create Hashlist -->
        <template v-if="inputs.activeStep == StepIndex.Hashlist">
          <div class="form-control">
            <HashlistInputs
              v-model:hashes="inputs.hashes"
              v-model:hashType="inputs.hashType"
              v-model:hashlistName="inputs.hashlistName"
            />

            <div class="mt-8 flex justify-between">
              <div class="flex justify-start">
                <button
                  class="link"
                  @click="saveOrGetHashlist"
                  v-if="hashlistStepValidationError == null"
                >
                  Save hashlist and finish
                </button>
              </div>
              <div class="card-actions justify-end">
                <button class="btn btn-ghost" @click="inputs.activeStep--">Previous</button>
                <div class="tooltip" :data-tip="hashlistStepValidationError">
                  <button
                    class="btn btn-primary"
                    @click="inputs.activeStep++"
                    :disabled="hashlistStepValidationError != null"
                  >
                    Next
                  </button>
                </div>
              </div>
            </div>
          </div>
        </template>

        <!-- Attack settings -->
        <template v-if="inputs.activeStep == StepIndex.Attack">
          <div class="join self-center">
            <input
              type="radio"
              name="options"
              :data-title="attackMode.name"
              class="btn btn-neutral join-item"
              :key="attackMode.value"
              :value="attackMode.value"
              v-model="inputs.attackMode"
              :aria-label="attackMode.name"
              v-for="attackMode in attackModes"
            />
          </div>

          <div class="my-2"></div>

          <!-- Wordlist -->
          <div v-if="inputs.attackMode === 0">
            <WordlistSelect
              label-text="Select Wordlist"
              :list="allWordlists?.wordlists ?? []"
              v-model="inputs.selectedWordlists"
              :limit="1"
            />
            <hr class="my-4" />
            <WordlistSelect
              label-text="Select Rule File(s)"
              :list="allRulefiles?.rulefiles ?? []"
              v-model="inputs.selectedRulefiles"
              :limit="Infinity"
            />
          </div>

          <!-- Combinator -->
          <div v-if="inputs.attackMode === 1">
            <WordlistSelect
              label-text="Select Left Wordlist"
              :list="allWordlists?.wordlists ?? []"
              v-model="inputs.combinatorLeft"
              :limit="1"
            />
            <hr class="my-4" />
            <WordlistSelect
              label-text="Select Right Wordlist"
              :list="allWordlists?.wordlists ?? []"
              v-model="inputs.combinatorRight"
              :limit="1"
            />
          </div>

          <!-- Brute-force/Mask -->
          <div v-if="inputs.attackMode === 3">
            <MaskInput v-model="inputs.mask" />
          </div>

          <!-- Wordlist + Mask -->
          <div v-if="inputs.attackMode === 6">
            <WordlistSelect
              label-text="Select Wordlist"
              :list="allWordlists?.wordlists ?? []"
              v-model="inputs.selectedWordlists"
              :limit="1"
            />
            <hr class="my-4" />
            <MaskInput v-model="inputs.mask" />
          </div>

          <!-- Mask + Wordlist -->
          <div v-if="inputs.attackMode === 7">
            <MaskInput v-model="inputs.mask" />
            <hr class="my-4" />
            <WordlistSelect
              label-text="Select Wordlist"
              :list="allWordlists?.wordlists ?? []"
              v-model="inputs.selectedWordlists"
              :limit="1"
            />
          </div>

          <hr class="my-4" />

          <label class="label font-bold">Additional Options</label>
          <div class="pl-3">
            <label class="label cursor-pointer justify-start">
              <input
                type="checkbox"
                v-model="inputs.isDistributed"
                class="checkbox-primary checkbox checkbox-xs"
              />
              <span><span class="label-text ml-4 font-bold">Distribute attack?</span></span>
            </label>
            <label class="label cursor-pointer justify-start">
              <input
                type="checkbox"
                v-model="inputs.enableLoopback"
                class="checkbox-primary checkbox checkbox-xs"
              />
              <span><span class="label-text ml-4 font-bold">Loopback?</span> (--loopback)</span>
            </label>
            <label class="label cursor-pointer justify-start">
              <input
                type="checkbox"
                v-model="inputs.optimizedKernels"
                class="checkbox-primary checkbox checkbox-xs"
              />
              <span><span class="label-text ml-4 font-bold">Optimized Kernels?</span> (-O)</span>
            </label>
            <label class="label cursor-pointer justify-start">
              <input
                type="checkbox"
                v-model="inputs.slowCandidates"
                class="checkbox-primary checkbox checkbox-xs"
              />
              <span><span class="label-text ml-4 font-bold">Slow Candidates?</span> (-S)</span>
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
        <template v-if="inputs.activeStep == StepIndex.Review">
          <table class="first-col-bold table w-full">
            <tbody>
              <div v-if="props.firstStep == 0">
                <tr>
                  <td>Project Name</td>
                  <td>{{ inputs.projectName }}</td>
                </tr>
                <tr>
                  <td>Project Description</td>
                  <td>{{ inputs.projectDesc }}</td>
                </tr>
              </div>
              <div v-if="props.firstStep <= 1">
                <tr>
                  <td>Hashlist Name</td>
                  <td>{{ inputs.hashlistName }}</td>
                </tr>
                <tr>
                  <td>Hashlist Type</td>
                  <td>{{ selectedHashType?.id }} - {{ selectedHashType?.name }}</td>
                </tr>
                <tr>
                  <td>Number of Hashes</td>
                  <td>{{ hashesArr.length }}</td>
                </tr>
              </div>
            </tbody>
          </table>

          <AttackConfigDetails
            :hashcatParams="computedHashcatParams"
            :is-distributed="inputs.isDistributed"
          ></AttackConfigDetails>

          <div class="mt-8 flex justify-between">
            <div class="flex justify-start">
              <button class="link" @click="saveUptoAttack">Save attack and finish</button>
            </div>

            <div class="card-actions justify-end">
              <button class="btn btn-ghost" @click="inputs.activeStep--">Previous</button>
              <button class="btn btn-success" @click="saveAndStartAttack">Start Attack</button>
            </div>
          </div>
        </template>
      </div>
    </div>
  </div>
</template>

<style scoped>
table.first-col-bold tr > td:first-of-type {
  font-weight: bold;
}
</style>
