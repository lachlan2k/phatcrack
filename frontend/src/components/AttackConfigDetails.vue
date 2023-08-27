<script setup lang="ts">
import { computed } from 'vue'

import type { AttackWithJobsDTO } from '@/api/types'
import { modeHasMask } from '@/util/hashcat'
import { useListfilesStore } from '@/stores/listfiles'

const props = defineProps<{
  attack: AttackWithJobsDTO
}>()

const listfileStore = useListfilesStore()

listfileStore.load()
function listfileName(id: string): string {
    return listfileStore.byId(id)?.name ?? 'Unknown'
}

const extraOptionsStr = computed(() => {
    const params = props.attack.hashcat_params
    const arr = []

    if (params.optimized_kernels) {
        arr.push('Optimized kernels')
    }

    if (params.slow_candidates) {
        arr.push('Slow candidates')
    }

    if (modeHasMask(params.attack_mode) && params.mask_increment) {
        if (params.mask_increment_max > 0) {
            arr.push(`Mask increment (min ${params.mask_increment_min}, max ${params.mask_increment_max})`)
        } else {
            arr.push('Mask increment')
        }
    }

    return arr.join(', ')
})
</script>

<template>
    <table class="compact-table table w-full">
      <thead>
        <tr>
          <th>Option</th>
          <th>Value</th>
        </tr>
      </thead>
      <tbody>
        <tr v-if="attack.hashcat_params.wordlist_filenames.length > 0">
          <td><strong>Wordlist</strong></td>
          <td>{{ attack.hashcat_params.wordlist_filenames.map(id => listfileName(id)).join(', ') }}</td>
        </tr>

        <tr v-if="attack.hashcat_params.rules_filenames.length > 0">
            <td><strong>Rules</strong></td>
            <td>{{ attack.hashcat_params.rules_filenames.map(id => listfileName(id)).join(', ') }}</td>
        </tr>

        <tr v-if="attack.hashcat_params.mask != ''">
            <td><strong>Mask</strong></td>
            <td>{{ attack.hashcat_params.mask }}</td>
        </tr>

        <tr v-if="extraOptionsStr != ''">
            <td><strong>Extra Options</strong></td>
            <td>{{ extraOptionsStr }}</td>
        </tr>
      </tbody>
    </table>
</template>