<!--
   * Author : see AUTHORS
   * Licence: MIT, see LICENSE
-->

<template>
  <span class="json-boolean">{{ jsonValue }}</span>
</template>

<script>
export default {
    name: 'JsonBoolean',
    props: {
        jsonValue: Boolean
    }
};
</script>

<style lang="scss">
.json-boolean {
    color: #fc1e70;
}
</style>
