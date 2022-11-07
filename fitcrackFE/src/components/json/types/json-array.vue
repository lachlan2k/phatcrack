<!--
   * Author : see AUTHORS
   * Licence: MIT, see LICENSE
-->

<template>
  <span>
    <v-icon
      v-if="!keyName"
      @click.stop="toggle"
    >
      {{ value ? 'arrow_drop_down' : 'arrow_drop_up' }}
    </v-icon>
    <span>[</span>
    <template v-if="jsonValue.length">
      <json-box
        v-for="(val, index) in jsonValue"
        v-show="value"
        :key="index"
        :value="val"
      />
      <span
        v-show="!value"
        class="node-ellipsis"
      >...</span>
    </template>
    <span>]</span>
  </span>
</template>

<script>
export default {
    name: 'JsonArray',
    props: {
        jsonValue: Array,
        keyName: String,
        value: Boolean
    },
    methods: {
        toggle() {
            this.$emit('input', !this.value);
        }
    }
};
</script>

<style>

</style>
