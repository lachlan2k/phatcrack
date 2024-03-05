import { ref, computed } from 'vue'

export function useHashesInput() {
  const hashesInput = ref('')
  const hashesArr = computed(() => {
    return hashesInput.value
      .trim()
      .split(/\n+/)
      .filter((x) => !!x)
      .map((x) => x.trim())
  })
  return { hashesInput, hashesArr }
}
