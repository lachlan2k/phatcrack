<script setup lang="ts">
import { reactive, watch } from 'vue'
import HashlistInputs from './Wizard/HashlistInputs.vue'
import { useApi } from '@/composables/useApi'
import { getHashlist } from '@/api/project'

const props = defineProps<{
  hashlistId: string
}>()

const inputs = reactive({
  hashlistName: '',
  hashType: '0',
  hashes: ''
})

const { data } = useApi(() => getHashlist(props.hashlistId))

watch(data, (newData) => {
  if (newData == null) {
    return
  }

  inputs.hashlistName = newData.name
  inputs.hashType = newData.hash_type.toString()
  inputs.hashes = newData.hashes.map((x) => x.input_hash).join('\n')
})
</script>

<template>
  <HashlistInputs
    :includeSaveButton="true"
    v-model:hashlistName="inputs.hashlistName"
    v-model:hashType="inputs.hashType"
    v-model:hashes="inputs.hashes"
  />
</template>
