<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { useToast } from 'vue-toastification'

import { bytesToReadable } from '@/util/units'
import { uploadListfile } from '@/api/listfiles'
import { useListfilesStore } from '@/stores/listfiles'
import type { AxiosProgressEvent } from 'axios'

enum FileType {
  Wordlist = 'Wordlist',
  Rulefile = 'Rulefile'
}

const MaxSizeInBytes = 10 * 1000 ** 3 // 10GB
const MaxSizeForAutoLineCount = 500 * 1000 ** 2 // 500MB

const props = defineProps<{
  fileType: FileType | null
}>()

const fileInputEl = ref<HTMLInputElement | null>(null)

const fileName = ref('')
const lineCount = ref(0)
const fileType = ref(props.fileType ?? FileType.Wordlist)
const fileToUpload = ref<File | null>(null)

const isLoading = ref(false)
const progress = ref<AxiosProgressEvent | null>(null)

const validationError = computed(() => {
  if (fileToUpload.value == null) {
    return 'Please select a file'
  }

  if (fileToUpload.value.size > MaxSizeForAutoLineCount && lineCount.value == 0) {
    return 'Please set the line count'
  }

  return null
})

const requiresLineCountSpecified = computed(() => {
  if (fileToUpload.value == null) {
    return false
  }

  return fileToUpload.value.size > MaxSizeForAutoLineCount
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

const listfilesStore = useListfilesStore()

async function onSubmit(event: Event) {
  event.preventDefault()
  if (fileToUpload.value == null) {
    return
  }

  const formData = new FormData()

  formData.append('file-type', props.fileType ?? fileType.value)
  formData.append('file-line-count', lineCount.value.toString())
  formData.append('file', fileToUpload.value)

  try {
    isLoading.value = true
    const uploadedFile = await uploadListfile(
      formData,
      (newProgress: AxiosProgressEvent) => (progress.value = newProgress)
    )
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
    if (e instanceof Error) {
      toast.error('Failed to upload file: ' + e.message)
    } else {
      toast.error('Failed to upload file')
    }
  } finally {
    isLoading.value = false
  }
}
</script>

<template>
  <h3 class="text-lg font-bold">Upload a {{ props.fileType == null ? 'File' : props.fileType }}</h3>
  <div class="form-control mt-1">
    <label class="label">
      <span class="label-text">Name</span>
    </label>
    <input
      type="text"
      class="input-bordered input"
      v-model="fileName"
      :placeholder="fileType == FileType.Rulefile ? 'best64.rule' : 'rockyou.txt'"
    />
  </div>

  <div class="form-control mt-1" v-if="requiresLineCountSpecified">
    <label class="label">
      <span class="label-text">Number of lines</span>
    </label>
    <input type="number" class="input-bordered input" v-model="lineCount" />
    <label class="label" v-if="lineCount == 0">
      <span class="label-text text-error"
        >Files larger {{ bytesToReadable(MaxSizeForAutoLineCount) }} require a line count</span
      >
    </label>
  </div>

  <div class="form-control mt-1" v-if="props.fileType == null">
    <label class="label">
      <span class="label-text">File type</span>
    </label>
    <select class="select-bordered select" v-model="fileType">
      <option value="Wordlist">Wordlist</option>
      <option value="Rulefile">Rulefile</option>
    </select>
  </div>

  <div class="form-control mt-1">
    <label class="label">
      <span class="label-text">Pick a file (max {{ bytesToReadable(MaxSizeInBytes) }})</span>
    </label>
    <input
      type="file"
      ref="fileInputEl"
      @change="onFileSelect"
      class="file-input-bordered file-input-ghost file-input"
      name="file"
    />
  </div>
  <div v-if="isLoading && progress != null && progress.total != null">
    <progress
      class="progress-primary progress w-full"
      :value="(progress.loaded / progress.total) * 100"
      max="100"
    ></progress>
  </div>

  <div class="form-control mt-3">
    <span class="tooltip" :data-tip="validationError">
      <button
        @click="onSubmit"
        :disabled="validationError != null || isLoading"
        class="btn-primary btn w-full"
      >
        <span class="loading loading-spinner loading-md" v-if="isLoading"></span>
        {{ buttonText }}
      </button>
    </span>
  </div>
</template>
