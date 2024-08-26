<script setup lang="ts">
import { computed } from 'vue'

const props = defineProps<{
  modelValue: {
    [key: string]: boolean
  }
}>()

const entries = computed(() => Object.entries(props.modelValue))

const emit = defineEmits(['update:modelValue'])

function toggleValue(key: string) {
  const newValue = {
    ...props.modelValue,
    [key]: !props.modelValue[key]
  }

  emit('update:modelValue', newValue)
}
</script>

<template>
  <div>
    <div v-for="[key, value] in entries" :key="key">
      <label
        class="label cursor-pointer"
        @click="
          (e) => {
            e.preventDefault()
            toggleValue(key)
          }
        "
      >
        <span class="label-text">{{ key }}</span>
        <input type="checkbox" :checked="value" class="checkbox" />
      </label>
    </div>
  </div>
</template>
