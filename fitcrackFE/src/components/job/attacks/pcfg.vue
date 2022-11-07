<!--
   * Author : see AUTHORS
   * Licence: MIT, see LICENSE
-->

<template>
  <div>
    <v-card-title>
      <span>Select PCFG grammar<span class="required primary--text"> *</span></span>
    </v-card-title>
    <pcfg-selector
      v-model="pcfg"
    />

    <v-divider />
    <v-card-title>
      <span>Edit keyspace limit</span>
    </v-card-title>
    <v-text-field
      v-model.number="keyspaceLimit"
      outlined
      single-line
      required
      type="number"
      suffix="passwords"
      :max="pcfg.length > 0 ? pcfg[0].keyspace : 0"
      @blur="checkValid"
    />
    <v-divider />
    <v-card-title>
      <span>Select rule file</span>
    </v-card-title>
    <rules-selector
      v-model="rules"
    />

    <v-divider />
    <v-checkbox
       v-model="optimized"
       label="Use optimized computing kernels (limits password length)"
    />
  </div>
</template>

<script>
  import pcfgSelector from '@/components/selector/pcfgSelector.vue'
  import ruleSelector from '@/components/selector/rulesSelector.vue'

  import {mapTwoWayState} from 'spyfu-vuex-helpers'
  import {twoWayMap} from '@/store'

  export default {
    name: "Pcfg",
    components: {
      'pcfg-selector': pcfgSelector,
      'rules-selector': ruleSelector
    },
    data: function () {
      return {
        attackId: 9,
        attackName: 'pcfg',
      }
    },
    watch: {
      pcfg (val) {
        if (val.length > 0) {
          this.keyspaceLimit = val[0].keyspace
        }
      }
    },
    computed: mapTwoWayState('jobForm', twoWayMap(['rules', 'pcfg', 'keyspaceLimit', 'optimized'])),
    methods: {
      checkValid: function () {
        if(this.pcfg.length > 0 && this.keyspaceLimit > this.pcfg[0].keyspace) {
          this.keyspaceLimit = this.pcfg[0].keyspace
        }
      }
    }
  }
</script>

<style scoped>
  .required {
    font-weight: bold;
  }
</style>
