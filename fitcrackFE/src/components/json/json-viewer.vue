<!--
   * Author : see AUTHORS
   * Licence: MIT, see LICENSE
-->

<template>
  <div class="json-viewer">
    <div class="code-box">
      <json-box :value="value" />
    </div>
  </div>
</template>

<script>
import Vue from 'vue';
import JsonBox from '@/components/json/json-box.vue'
Vue.component(JsonBox.name, JsonBox);
export default {
    name: 'JsonViewer',
     components: {
        JsonBox
    },
    props: {
        value: [Object, Array, String, Number, Boolean]
    }
};
</script>

<style lang="scss">
.json-viewer {
    box-sizing: border-box;
    position: relative;

    .code-box {
        padding: 20px;
    }

}
</style>
