<!--
   * Author : see AUTHORS
   * Licence: MIT, see LICENSE
-->

<template>
  <div>
    <v-row>
      <v-col
        cols="6"
        class="border"
      >
        <v-card-title>
          <span>Select left dictionary<span class="required primary--text"> *</span></span>
        </v-card-title>
      </v-col>
      <v-col cols="6">
        <v-card-title>
          <span>Select right dictionary<span class="required primary--text"> *</span></span>
        </v-card-title>
      </v-col>
      <v-col
        cols="6"
        class="border"
      >
        <dict-selector
          v-model="leftDicts"
          select-all
          @input="checkValid"
        />
      </v-col>
      <v-col cols="6">
        <dict-selector
          v-model="rightDicts"
          select-all
          @input="checkValid"
        />
      </v-col>

      <v-col
        cols="6"
        class="border"
      >
        <v-card-title>
          <span>Type left rule</span>
        </v-card-title>
      </v-col>
      <v-col cols="6">
        <v-card-title>
          <span>Type right rule</span>
        </v-card-title>
      </v-col>

      <v-col
        cols="6"
        class="border px-2"
      >
        <v-text-field
          v-model="ruleLeft"
          outlined
          single-line
          placeholder="Rule"
          @input="checkValid"
        />
      </v-col>
      <v-col
        cols="6"
        class="px-2"
      >
        <v-text-field
          v-model="ruleRight"
          outlined
          single-line
          placeholder="Rule"
          @input="checkValid"
        />
      </v-col>
    </v-row>

    <v-divider />
    <v-checkbox
       v-model="optimized"
       label="Use optimized computing kernels (limits password length)"
    />
  </div>
</template>

<script>
  import dictSelector from '@/components/selector/dictionarySelector.vue'

  import {mapTwoWayState} from 'spyfu-vuex-helpers'
  import {twoWayMap} from '@/store'

  export default {
    name: "Combinator",
    components: {
      'dict-selector': dictSelector
    },
    computed: mapTwoWayState('jobForm', twoWayMap(['leftDicts', 'ruleLeft', 'rightDicts', 'ruleRight', 'optimized'])),
    methods: {
      checkValid: function () {
        if (this.leftDicts.length > 0 && this.rightDicts.length > 0) {
          return true
        }
        return false
      }
    }
  }
</script>

<style scoped>
  .border {
    border-right: 2px solid rgba(0,0,0,.12)
  }

  .required {
    font-weight: bold;
    color: #000;
  }
</style>

