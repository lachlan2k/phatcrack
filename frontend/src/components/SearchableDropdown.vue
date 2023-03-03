<script setup lang="ts">
import { ref, computed } from 'vue'

interface OptionT {
  value: string
  text: string
}

const props = defineProps<{
  modelValue: string
  placeholderText: string
  options: OptionT[]
}>()

const emit = defineEmits(['update:modelValue'])

const inputText = ref('')

const filteredOptions = computed(() =>
  props.options.filter((x) => x.text.toLowerCase().includes(inputText.value.toLowerCase()))
)

const optionsVisible = ref(false)

function selectOption(option: OptionT) {
  optionsVisible.value = false
  emit('update:modelValue', option.value)
  inputText.value = option.text
}

function focus() {
  optionsVisible.value = true
  inputText.value = ''
}

function unfocus() {
  optionsVisible.value = false
}
</script>

<template>
  <div class="relative">
    <input
      type="text"
      class="input-bordered input w-full cursor-pointer focus:outline-none"
      :placeholder="props.placeholderText"
      v-model="inputText"
      @focus="focus"
      @blur="unfocus"
    />
    <div
      v-if="optionsVisible"
      class="floating-dropdown-content absolute w-full border-solid border-black shadow-md"
    >
      <div
        :key="option.value"
        v-for="option in filteredOptions"
        class="dropdown-content-option hover mx-1 my-1 cursor-pointer px-2 py-1"
        @mousedown="selectOption(option)"
      >
        {{ option.text }}
      </div>
    </div>
  </div>
</template>

<style scoped>
.floating-dropdown-content {
  background: white;
  max-height: 400px;
  overflow-y: scroll;
  overflow-x: hidden;
}

.dropdown-content-option {
  width: 100%;
}

.dropdown-content-option:hover {
  width: 100%;
  background: #eee;
}
</style>
