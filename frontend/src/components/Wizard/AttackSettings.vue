<script setup lang="ts">
import { storeToRefs } from 'pinia'
import { computed, watch } from 'vue'

import MaskInput from '@/components/Wizard/MaskInput.vue'
import WordlistSelect from '@/components/Wizard/ListSelect.vue'

import { useListfilesStore } from '@/stores/listfiles'

import type { AttackMode } from '@/util/hashcat'
import { attackModes } from '@/util/hashcat'

export interface AttackSettingsT {
  attackMode: AttackMode

  selectedWordlists: string[]
  selectedRulefiles: string[]

  mask: string
  maskIncrement: boolean

  combinatorLeft: string[]
  combinatorRight: string[]

  optimizedKernels: boolean
  slowCandidates: boolean
  enableLoopback: boolean

  isDistributed: boolean
}

const props = defineProps<{
  modelValue: AttackSettingsT
}>()

const emit = defineEmits(['update:modelValue'])

const attackSettings = computed({
  get: () => props.modelValue,
  set: newVal => emit('update:modelValue', newVal)
})

const listfileStore = useListfilesStore()
listfileStore.load(true)
const { wordlists, rulefiles } = storeToRefs(listfileStore)

watch(
  () => attackSettings.value.combinatorLeft,
  newLeft => (attackSettings.value.selectedWordlists = [...newLeft, ...attackSettings.value.combinatorRight])
)
watch(
  () => attackSettings.value.combinatorRight,
  newRight => (attackSettings.value.selectedWordlists = [...attackSettings.value.combinatorLeft, ...newRight])
)
</script>

<template v-if="attackSettings.activeStep == StepIndex.Attack">
  <div class="join self-center">
    <input
      type="radio"
      name="options"
      :data-title="attackMode.name"
      class="btn btn-neutral join-item"
      :key="attackMode.value"
      :value="attackMode.value"
      v-model="attackSettings.attackMode"
      :aria-label="attackMode.name"
      v-for="attackMode in attackModes"
    />
  </div>

  <div class="my-2"></div>

  <!-- Wordlist -->
  <div v-if="attackSettings.attackMode === 0">
    <WordlistSelect label-text="Select Wordlist" :list="wordlists" v-model="attackSettings.selectedWordlists" :limit="1" />
    <hr class="my-4" />
    <WordlistSelect label-text="Select Rule File(s)" :list="rulefiles" v-model="attackSettings.selectedRulefiles" :limit="Infinity" />
  </div>

  <!-- Combinator -->
  <div v-if="attackSettings.attackMode === 1">
    <WordlistSelect label-text="Select Left Wordlist" :list="wordlists" v-model="attackSettings.combinatorLeft" :limit="1" />
    <hr class="my-4" />
    <WordlistSelect label-text="Select Right Wordlist" :list="wordlists" v-model="attackSettings.combinatorRight" :limit="1" />
  </div>

  <!-- Brute-force/Mask -->
  <div v-if="attackSettings.attackMode === 3">
    <MaskInput v-model="attackSettings.mask" />
    <label class="label cursor-pointer justify-start">
      <input type="checkbox" v-model="attackSettings.maskIncrement" class="checkbox-primary checkbox checkbox-xs" />
      <span><span class="label-text ml-4 font-bold">Mask increment?</span></span>
    </label>
  </div>

  <!-- Wordlist + Mask -->
  <div v-if="attackSettings.attackMode === 6">
    <WordlistSelect label-text="Select Wordlist" :list="wordlists" v-model="attackSettings.selectedWordlists" :limit="1" />
    <hr class="my-4" />
    <MaskInput v-model="attackSettings.mask" />
    <label class="label cursor-pointer justify-start">
      <input type="checkbox" v-model="attackSettings.maskIncrement" class="checkbox-primary checkbox checkbox-xs" />
      <span><span class="label-text ml-4 font-bold">Mask increment?</span></span>
    </label>
  </div>

  <!-- Mask + Wordlist -->
  <div v-if="attackSettings.attackMode === 7">
    <MaskInput v-model="attackSettings.mask" />
    <hr class="my-4" />
    <WordlistSelect label-text="Select Wordlist" :list="wordlists" v-model="attackSettings.selectedWordlists" :limit="1" />
    <label class="label cursor-pointer justify-start">
      <input type="checkbox" v-model="attackSettings.maskIncrement" class="checkbox-primary checkbox checkbox-xs" />
      <span><span class="label-text ml-4 font-bold">Mask increment?</span></span>
    </label>
  </div>

  <hr class="my-4" />

  <label class="label font-bold">Additional Options</label>
  <div>
    <label class="label cursor-pointer justify-start">
      <input type="checkbox" v-model="attackSettings.isDistributed" class="checkbox-primary checkbox checkbox-xs" />
      <span><span class="label-text ml-4 font-bold">Distribute attack?</span></span>
    </label>
    <label class="label cursor-pointer justify-start">
      <input type="checkbox" v-model="attackSettings.enableLoopback" class="checkbox-primary checkbox checkbox-xs" />
      <span><span class="label-text ml-4 font-bold">Loopback?</span> (--loopback)</span>
    </label>
    <label class="label cursor-pointer justify-start">
      <input type="checkbox" v-model="attackSettings.optimizedKernels" class="checkbox-primary checkbox checkbox-xs" />
      <span><span class="label-text ml-4 font-bold">Optimized Kernels?</span> (-O)</span>
    </label>
    <label class="label cursor-pointer justify-start">
      <input type="checkbox" v-model="attackSettings.slowCandidates" class="checkbox-primary checkbox checkbox-xs" />
      <span><span class="label-text ml-4 font-bold">Slow Candidates?</span> (-S)</span>
    </label>
  </div>
</template>
