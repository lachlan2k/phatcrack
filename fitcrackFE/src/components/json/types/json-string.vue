<!--
   * Author : see AUTHORS
   * Licence: MIT, see LICENSE
-->

<template>
  <span class="json-string">"{{ jsonValue }}"</span>
</template>

<script>
export default {
    name: 'JsonString',
    props: {
        jsonValue: String
    }
};
</script>

<style lang="scss">
.json-string {
    color: #42b983;
}
</style>
