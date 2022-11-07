<script setup lang="ts">
const version = 'v1.0.0'
</script>

<template>
<div class="version">
    {{ version }}
</div>
</template>

<style scoped>
.version {
    position: fixed;
    bottom: 1.5em;
    left: 0;
    width: 100%;
    font-size: 0.8em;
    font-weight: bold;
    color: #888;
    text-align: center;
}
</style>