<script setup lang="ts">
import { ref, computed, watch } from 'vue'

interface OptionT {
    value: string
    text: string
}

const props = defineProps<{
    modelValue: string
    options: OptionT[]
}>()

const emit = defineEmits(['update:modelValue'])

// const selectedIndex = ref(options.findIndex(option => option.value === modelValue))

const placeholderText = 'Choose hash type'
const inputText = ref('')

const filteredOptions = computed(() => props.options.filter(x => x.text.toLowerCase().includes(inputText.value.toLowerCase())))

const optionsVisible = ref(false)

function selectIndex(index: number) {
    console.log('settig index', index)
    // selectedIndex.value = index
    optionsVisible.value = false
    inputText.value = props.options[index].text
    emit('update:modelValue', props.options[index].value)
}

function focus() {
    optionsVisible.value = true
    inputText.value = ''
}

function unfocus() {
    // if (selectedIndex.value > -1) {
        // inputText.value = options[selectedIndex.value].text
    // }
    optionsVisible.value = false
}

</script>

<template>
    <div class="relative">
        <input type="text" class="cursor-pointer input-bordered input w-full" :placeholder="placeholderText" v-model="inputText" @focus="focus" @blur="unfocus">
        <div v-if="optionsVisible" class="floating-dropdown-content absolute w-full border-solid border-black">
            <div v-for="option, optionIndex in filteredOptions" class="dropdown-content-option mx-1 my-1 px-2 py-1 cursor-pointer hover" @mousedown="selectIndex(optionIndex)">{{ option.text }}</div>
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