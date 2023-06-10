<script setup lang="ts">
import SearchableDropdown from '@/components/SearchableDropdown.vue'
import { computed, watch } from 'vue'
import { storeToRefs } from 'pinia'
import { useWizardHashDetect } from '@/composables/useWizardHashDetect'
import { useResourcesStore } from '@/stores/resources'

const props = defineProps<{
  hashlistName: string
  hashType: string
  hashes: string
  includeSaveButton?: boolean
}>()

const emit = defineEmits(['update:hashlistName', 'update:hashes', 'update:hashType', 'savePressed'])

const resourcesStore = useResourcesStore()
const { hashTypes: allHashTypes } = storeToRefs(resourcesStore)
resourcesStore.loadHashTypes()

const hashes = computed({
  get: () => props.hashes,
  set: (value: string) => emit('update:hashes', value)
})

const hashlistName = computed({
  get: () => props.hashlistName,
  set: (value: string) => emit('update:hashlistName', value)
})

const hashType = computed({
  get: () => props.hashType,
  set: (value: string) => emit('update:hashType', value)
})

const hashesArr = computed(() => {
  return hashes.value.split(/\s+/).filter((x) => !!x)
})

const {
  detectButtonClass,
  detectButtonClick,
  detectButtonText,
  suggestedHashTypes,
  isLoadingSuggestions
} = useWizardHashDetect(hashesArr)

watch(suggestedHashTypes, (newHashTypes) => {
  const types = newHashTypes?.possible_types
  if (!types || types.length == 0) {
    return
  }
  hashType.value = types.sort()[0].toString()
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
</script>

<template>
  <label class="label font-bold">
    <span class="label-text">Hashlist Name</span>
  </label>
  <input
    type="text"
    placeholder="Dumped Admin NTLM Hashes"
    v-model="hashlistName"
    class="input-bordered input w-full max-w-xs"
  />
  <hr class="my-4" />
  <label class="label font-bold">
    <span class="label-text">Hash Type ({{ filteredHashTypes.length }} options)</span>
  </label>

  <div class="flex justify-between">
    <SearchableDropdown
      class="flex-grow"
      v-model="hashType"
      :options="hashTypeOptionsToShow"
      placeholder-text="Search for a hashtype..."
    />
    <button
      class="btn ml-1"
      :class="detectButtonClass"
      :disabled="isLoadingSuggestions || hashesArr.length == 0"
      @click="detectButtonClick"
    >
      {{ detectButtonText }}
    </button>
  </div>
  <label class="label mt-2 font-bold">
    <span class="label-text">Hashes (one per line)</span>
  </label>
  <textarea
    placeholder="Hashes"
    class="hashes-input textarea-bordered textarea w-full font-mono focus:outline-none"
    rows="8"
    v-model="hashes"
  ></textarea>
  <div v-if="props.includeSaveButton" class="mt-4 flex justify-end">
    <button class="btn-primary btn" @click="emit('savePressed')">Save</button>
  </div>
</template>

<style scoped>
textarea.hashes-input {
  background-image: linear-gradient(to bottom, rgba(87, 87, 87, 0.05) 50%, transparent 50%);
  background-repeat: repeat-y;

  background-size: 100% 50px;

  line-height: 25px;
  font-size: 15px;

  padding: 0;
  padding-left: 3px;

  white-space: pre;
  resize: none;

  background-attachment: local;
}
</style>
