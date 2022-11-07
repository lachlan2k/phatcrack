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
          <span>Select dictionary<span class="required primary--text"> *</span></span>
        </v-card-title>
      </v-col>
      <v-col cols="6">
        <v-card-title>
          <span>Type mask<span class="required primary--text"> *</span></span>
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
        <v-card-title>
          <span>Type rule</span>
        </v-card-title>
        <v-text-field
          v-model="ruleLeft"
          outlined
          single-line
          placeholder="Rule"
          @input="checkValid"
        />
      </v-col>
      <v-col cols="6">
        <mask-single
          v-model="hybridMask"
          non-removable
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
  import maskSingle from '@/components/mask/maskSingle.vue'

  import {mapTwoWayState} from 'spyfu-vuex-helpers'
  import {twoWayMap} from '@/store'

  export default {
    name: "HybridMaskWordlist",
    components: {
      'mask-single': maskSingle,
      'dict-selector': dictSelector
    },
    data: function () {
      return {
        attackId: 6,
        attackName: 'Hybrid wordlist+mask',
        maskRules: [
          v => /^(\?[ludhHsab]|[ -~])+$/.test(v) || 'Not valid mask'
        ]
      }
    },
    computed: mapTwoWayState('jobForm', twoWayMap(['ruleLeft', 'leftDicts', 'ruleRight', 'hybridMask', 'optimized'])),
    methods: {
      checkValid: function () {
        if (this.hybridMask !== '' && this.leftDicts.length > 0) {
          /*
          this.$emit('input', {
            'attack_mode': this.attackId,
            'attack_name': this.attackName,
            'rule_left': this.ruleLeft,
            'rule_right': this.ruleRight,
            'left_dictionaries': this.dictionaries,
            'mask': this.mask
          })
          */
          return true
        }
        return false
      },
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

  .maskBtn {
    min-width: 0px;
    max-width: 40px;
    text-transform: none;
  }

  .maskBtns {
    text-align: center;
  }

</style>

<style>
  .selected {
    background: rgba(37, 157, 173, 0.85) !important;
    color: white;
  }

  .selected a {
    color: white;
  }

  .clickable {
    cursor: pointer;
  }
</style>
