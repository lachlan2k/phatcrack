<!--
   * Author : see AUTHORS
   * Licence: MIT, see LICENSE
-->

<template>
  <div>
    <v-card-title class="pb-0 mb-2">
      <span>Add masks<span class="required primary--text"> *</span></span>
    </v-card-title>

    <v-row class="mb-4 px-4">
      <v-btn
        color="primary"
        text
        @click="loadMasksDialog = true"
      >
        <v-icon left>
          mdi-file-download
        </v-icon>
        Load masks
      </v-btn>
      <v-btn
        v-show="masks.length > 1"
        color="error"
        text
        @click="masks = ['']"
      >
        <v-icon left>
          mdi-restart
        </v-icon>
        Reset masks
      </v-btn>
      <v-spacer />
      <v-btn
        id="add-mask-button"
        v-show="masks.length < 50"
        color="success"
        @click="addMask()"
      >
        <v-icon left>
          mdi-plus
        </v-icon>
        Add mask
      </v-btn>
    </v-row>

    <div
      v-if="masks.length < 50"
      class="mask-editors"
    >
      <mask-single
        v-for="(mask, i) in masks"
        :id="`mask-${i}`"
        :key="i"
        :non-removable="masks.length == 1"
        :custom-charsets="charset"
        :value="mask"
        @input="e => updateMask({index: i, val: e})"
        @remove="deleteMask(i)"
      />
    </div>
    <v-alert
      v-else
      type="info"
    >
      Not showing {{ masks.length }} masks to maintain performance and conciseness.
    </v-alert>

    <v-row>
      <v-col
        cols="6"
        class="border"
      >
        <v-card-title>
          <span>Select charsets (max. 4)</span>
        </v-card-title>
        <charset-selector
          v-model="charset"
          select-all
        />
      </v-col>
      <v-col cols="6">
        <v-card-title>
          <span>Markov file</span>
        </v-card-title>
        <markov-selector
          v-model="markov"
          :markov-submode="submode"
          @input="checkValid"
        />

        <v-row>
          <v-col
            cols="6"
            class="pl-2"
          >
            <v-radio-group
              v-model="submode"
              column
              @change="checkValid"
            >
              <v-radio
                label="Markov disabled"
                color="primary"
                :value="0"
              />
              <v-radio
                label="2D Markov"
                color="primary"
                :value="1"
              />
              <v-radio
                label="3D Markov"
                color="primary"
                :value="2"
              />
            </v-radio-group>
          </v-col>
          <v-col
            cols="6"
            class="pr-2"
          >
            <v-text-field
              v-model="markovThresh"
              type="number"
              :disabled="submode === 0"
              label="Markov threshold"
              filled
              dense
              mask="########"
              @input="checkValid"
            />
          </v-col>
        </v-row>
      </v-col>
    </v-row>



    <v-dialog
      v-model="loadMasksDialog"
      max-width="600"
    >
      <v-card>
        <v-card-text>
          <maskFile-selector
            v-model="maskFile"
            @input="checkValid"
          />
        </v-card-text>
        <v-card-actions>
          <v-btn
            color="primary"
            :disabled="!maskFile"
            @click="loadMasksFromFile"
          >
            Load
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <v-divider />
    <v-checkbox
       v-model="optimized"
       label="Use optimized computing kernels (limits password length)"
    />
  </div>
</template>

<script>
  import maskSingle from '@/components/mask/maskSingle.vue'
  import CharsetSelector from '@/components/selector/charsetSelector.vue'
  import markovSelector from '@/components/selector/markovSelector.vue'
  import maskFileSelector from '@/components/selector/loadMaskFileSelector.vue'

  import {mapMutations} from 'vuex'
  import {mapTwoWayState} from 'spyfu-vuex-helpers'
  import {twoWayMap} from '@/store'

  export default {
    name: "FcMask",
    components: {
      CharsetSelector,
      'mask-single': maskSingle,
      'charset-selector': CharsetSelector,
      'markov-selector': markovSelector,
      'maskFile-selector': maskFileSelector
    },
    data: function () {
      return {
        loadMasksDialog: false,
        maskFile: null
      }
    },
    computed: {
      ...mapTwoWayState('jobForm', twoWayMap(['masks', 'submode', 'markovThresh', 'markov', 'charset', 'optimized']))
    },
    methods: {
      ...mapMutations('jobForm', ['addMask', 'deleteMask', 'updateMask', 'mergeMasks']),
      loadMasksFromFile: function () {
        this.axios.get(this.$serverAddr + '/masks/' + this.maskFile.id).then((response) => {
          var res = response.data.data.split("\n")
          this.mergeMasks(res)
          this.maskFile = null
          this.loadMasksDialog = false
          this.checkValid()
        })
      },
      checkValid: function () {
        if (this.masks[0] !== '') {
          /*
          this.$emit('input', {
            'attack_mode': this.attackId,
            'attack_name': this.attackName,
            'masks': this.masks,
            'attack_submode': this.markovSubmode,
            'markov_treshold': this.markovSubmode == 0 ? null : parseInt(this.markovTreshold),
            'markov': this.markovSubmode == 0 ? null : this.markov[0],
            'charset': this.charsets
          })
          */
          return true
        }
        return false
      }
    }
  }
</script>

<style scoped>
  .mask-editors {
    display: flex;
    flex-direction: column-reverse;
  }

  .masksContainer {
    overflow: auto;
    max-height: 500px;
  }
  .required {
    font-weight: bold;
    color: #000;
  }

  .maskBtnAction {
    margin-top: 15px;
  }

  .maskBtn {
    width: 40px;
    min-width: 0px;
    text-transform: none;
  }

  .border {
    border-right: 2px solid rgba(0,0,0,.12)
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

  .width96 {
    width: 96%;
  }
</style>
