<script setup lang="ts">
import { ref } from 'vue'
import { createProject } from '@/api/project'

const props = defineProps<{
  // Set to 0 for full wizard, 1 if project is already made, 2 if hashlist is already made...
  firstStep?: number
  existingProjectID?: string
  existingHashlistId?: string
}>()

const activeStep = ref(0)

const steps = [
  { name: 'Name Project' },
  { name: 'Add Hashlist' },
  { name: 'Configure Attack Settings' },
  { name: 'Review & Start Attack' }
].slice(props.firstStep ?? 0)

const projectName = ref('')
const projectDesc = ref('')

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
          <input type="text" placeholder="Hash Type" class="input-bordered input w-full max-w-xs" />
          <textarea placeholder="Hashes" class="textarea-bordered textarea max-w-xs"></textarea>

          <div class="mt-8 flex justify-between">
            <div class="flex justify-start">
              <button class="link" @click="saveUptoHashlist">Save hashlist and finish</button>
            </div>

            <div class="card-actions justify-end">
              <button class="btn-ghost btn" @click="activeStep--">Previous</button>
              <button class="btn-primary btn" @click="activeStep++">Next</button>
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
