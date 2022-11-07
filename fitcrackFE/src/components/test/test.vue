<!--
   * Author : see AUTHORS
   * Licence: MIT, see LICENSE
-->

<template>
  <v-container>
    <v-row>
      <v-col cols="12">
        <test-element>
          test
        </test-element>
      </v-col>
    </v-row>
  </v-container>
</template>

<script>
  import element from '@/components/logs/logsView.vue'
  export default {
    name: "Test",
    components: {
      'test-element': element
    },
    data() {
      return {
      }
    },
  }
</script>

<style scoped>

</style>
