import { computed, type Ref } from 'vue'
import { detectHashType } from '@/api/hashcat'
import { useApi } from '@/composables/useApi'

export function useWizardHashDetect(hashesArr: Ref<string[]>) {
  const {
    fetchData: fetchHashTypeSuggestions,
    isLoading: isLoadingSuggestions,
    data: suggestedHashTypes
  } = useApi(() => detectHashType(hashesArr.value[0]), { immediate: false, toastOnError: true })

  function detectButtonClick() {
    // Reset
    if (suggestedHashTypes.value != null) {
      suggestedHashTypes.value = null
      return
    }

    fetchHashTypeSuggestions()
  }

  const detectStatusText = computed(() => {
    if (suggestedHashTypes.value == null) {
      return ''
    }

    if (suggestedHashTypes.value.possible_types.length > 0) {
      return `Filtered down to ${suggestedHashTypes.value.possible_types.length} possible hash types`
    }

    return 'No suggestions found, check your hashes are valid'
  })

  const detectButtonClass = computed(() => {
    if (isLoadingSuggestions.value) {
      return 'btn-ghost'
    }

    if (suggestedHashTypes.value != null) {
      return ''
    }

    return 'btn-neutral'
  })

  const detectButtonText = computed(() => {
    if (isLoadingSuggestions.value) {
      return 'Loading suggestions...'
    }

    if (suggestedHashTypes.value != null) {
      return 'Reset Filter'
    }

    return 'Detect'
  })

  return {
    detectButtonClass,
    detectButtonClick,
    detectButtonText,
    detectStatusText,
    suggestedHashTypes,
    isLoadingSuggestions
  }
}
