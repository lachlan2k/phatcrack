<script setup lang="ts">
import { computed, ref, watch } from 'vue'

interface FileT {
  id: string
  name: string
  description: string
  filename: string
  size: number
  lines: number
}

const props = defineProps<{
  modelValue: string[]

  files?: FileT[]
  allowMultiple: boolean

  openButtonText?: string
  btnClass?: string
  modalTitle?: string
  modalDescription?: string
}>()

const emit = defineEmits(['update:modelValue'])

const selected = ref(props.modelValue)
watch(
  () => props.modelValue,
  (newModelVal) => {
    selected.value = newModelVal
  }
)

const okButtonText = computed(() => {
  if (!props.allowMultiple) {
    return 'Select file'
  }

  if (selected.value.length == 1) {
    return 'Select 1 file'
  }

  return `Select ${selected.value.length} files`
})

function addToSelected(id: string) {
  console.log('click')
  if (!selected.value.includes(id)) {
    selected.value = [...selected.value, id]
  } else {
    console.log(selected.value)
    selected.value = selected.value.filter((x) => x != id)
    console.log(selected.value)
  }
}

function onClose() {
  emit('update:modelValue', selected.value)
}
</script>

<template>
  <div>
    <label for="my-modal" class="btn" :class="props.btnClass">{{
      props.openButtonText || 'Select file'
    }}</label>

    <input type="checkbox" id="my-modal" class="modal-toggle" />
    <div class="modal">
      <div class="modal-box">
        <h3 class="text-lg font-bold">
          {{ props.modalTitle || (props.allowMultiple ? 'Select files' : 'Select a file') }}
        </h3>
        <p class="py-4" v-if="props.modalDescription != null">{{ props.modalDescription }}</p>
        <table class="table w-full">
          <tbody>
            <tr>
              <td>Select</td>
              <td>Name</td>
              <td>Number of lines</td>
            </tr>
            <tr
              v-for="file in props.files"
              :key="file.id"
              @click="addToSelected(file.id)"
              class="cursor-pointer"
            >
              <td>
                <input
                  type="checkbox"
                  class="checkbox-primary checkbox checkbox-xs align-middle"
                  :checked="selected.includes(file.id)"
                />
              </td>
              <td>{{ file.name }}</td>
              <td>{{ file.lines }}</td>
            </tr>
          </tbody>
        </table>
        <div class="modal-action">
          <label for="my-modal" class="btn btn-primary" @click="onClose">{{ okButtonText }}</label>
        </div>
      </div>
    </div>
  </div>
</template>
