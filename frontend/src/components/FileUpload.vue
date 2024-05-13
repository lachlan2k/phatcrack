<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { useToast } from 'vue-toastification'

import { bytesToReadable } from '@/util/units'
import { uploadListfile } from '@/api/listfiles'
import { useListfilesStore } from '@/stores/listfiles'
import type { AxiosProgressEvent } from 'axios'
import { useConfigStore } from '@/stores/config'
import { storeToRefs } from 'pinia'
import { useAuthStore } from '@/stores/auth'
import { useToastError } from '@/composables/useToastError'

type FileType = 'Wordlist' | 'Rulefile'

const configStore = useConfigStore()
const { config } = storeToRefs(configStore)
configStore.load()

const authStore = useAuthStore()
const { isAdmin } = storeToRefs(authStore)

const props = defineProps<{
  fileType: FileType | null
}>()

const fileInputEl = ref<HTMLInputElement | null>(null)

const fileName = ref('')
const lineCount = ref(0)
const fileType = ref(props.fileType ?? 'Wordlist')
const fileToUpload = ref<File | null>(null)

const isLoading = ref(false)
const progress = ref<AxiosProgressEvent | null>(null)

const validationError = computed(() => {
  const conf = config.value
  if (conf == null) {
    return 'Required config failed to load'
  }

  if (fileToUpload.value == null) {
    return 'Please select a file'
  }

  if (fileToUpload.value.size > conf.general.maximum_uploaded_file_line_scan_size && lineCount.value == 0) {
    return 'Please set the line count'
  }

  if (fileToUpload.value.size > conf.general.maximum_uploaded_file_size && !isAdmin) {
    return 'File is too large'
  }

  return null
})

const requiresLineCountSpecified = computed(() => {
  if (fileToUpload.value == null) {
    return false
  }

  return fileToUpload.value.size > (config.value?.general.maximum_uploaded_file_line_scan_size ?? 0)
})

watch(requiresLineCountSpecified, (doesRequire) => {
  if (doesRequire) {
    // Set it back to 0 to ask the server to calculate it
    lineCount.value = 0
  }
})

const buttonText = computed(() => {
  if (validationError.value != null) {
    return 'Upload'
  }

  const verb = isLoading.value ? 'Uploading' : 'Upload'

  return `${verb} ${fileName.value} (${bytesToReadable(fileToUpload.value!.size)})`
})

async function onFileSelect(event: Event) {
  const target = event.target as HTMLInputElement
  const lastFileName = fileToUpload.value?.name ?? ''

  fileToUpload.value = target.files?.[0] ?? null

  if (fileName.value == lastFileName) {
    fileName.value = fileToUpload.value?.name ?? ''
  }
}

const toast = useToast()
const { catcher } = useToastError()

const listfilesStore = useListfilesStore()

async function onSubmit(event: Event) {
  event.preventDefault()
  if (fileToUpload.value == null) {
    return
  }

  const formData = new FormData()

  formData.append('file-name', fileName.value)
  formData.append('file-type', props.fileType ?? fileType.value)
  formData.append('file-line-count', lineCount.value.toString())
  formData.append('file', fileToUpload.value)

  try {
    isLoading.value = true
    const uploadedFile = await uploadListfile(formData, (newProgress: AxiosProgressEvent) => (progress.value = newProgress))
    toast.success('Successfully uploaded file: ' + uploadedFile.name)
    listfilesStore.load(true)

    fileName.value = ''
    fileToUpload.value = null
    lineCount.value = 0
    progress.value = null

    if (fileInputEl.value != null) {
      fileInputEl.value.value = ''
    }
  } catch (e) {
    catcher(e)
  } finally {
    isLoading.value = false
  }
}
</script>

<template>
  <h3 class="text-lg font-bold">Upload a {{ props.fileType == null ? 'File' : props.fileType }}</h3>
  <div class="form-control mt-1">
    <label class="label font-bold">
      <span class="label-text">Name</span>
    </label>
    <input
      type="text"
      class="input input-bordered"
      v-model="fileName"
      :placeholder="fileType == 'Rulefile' ? 'best64.rule' : 'rockyou.txt'"
    />
  </div>

  <div class="form-control mt-1" v-if="requiresLineCountSpecified">
    <label class="label font-bold">
      <span class="label-text">Number of lines</span>
    </label>
    <input type="number" class="input input-bordered" v-model="lineCount" />
    <label class="label" v-if="lineCount == 0">
      <span class="label-text text-error"
        >Files larger {{ bytesToReadable(config?.general.maximum_uploaded_file_line_scan_size ?? 0) }} require a line count</span
      >
    </label>
  </div>

  <div class="form-control mt-1" v-if="props.fileType == null">
    <label class="label font-bold">
      <span class="label-text">File type</span>
    </label>
    <select class="select select-bordered" v-model="fileType">
      <option value="Wordlist">Wordlist</option>
      <option value="Rulefile">Rulefile</option>
    </select>
  </div>

  <div class="form-control mt-1">
    <label class="label font-bold">
      <span class="label-text" v-if="isAdmin">Pick a file</span>
      <span class="label-text" v-else>Pick a file (max {{ bytesToReadable(config?.general.maximum_uploaded_file_size ?? 0) }})</span>
    </label>
    <input type="file" ref="fileInputEl" @change="onFileSelect" class="file-input file-input-bordered file-input-ghost" name="file" />
  </div>
  <div v-if="isLoading && progress != null && progress.total != null">
    <progress class="progress progress-primary w-full" :value="(progress.loaded / progress.total) * 100" max="100"></progress>
  </div>

  <div class="form-control mt-3">
    <span class="tooltip" :data-tip="validationError">
      <button @click="onSubmit" :disabled="validationError != null || isLoading" class="btn btn-primary w-full">
        <span class="loading loading-spinner loading-md" v-if="isLoading"></span>
        {{ buttonText }}
      </button>
    </span>
  </div>
</template>
