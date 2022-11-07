<!--
   * Author : see AUTHORS
   * Licence: MIT, see LICENSE
-->

<template>
  <fc-tile title="Attack details">
    <v-data-table
      :headers="maskHeaders"
      :items="data.masks"
      hide-default-footer
    >
      <template v-slot:item.progress="{ item }">
        <v-progress-circular
          size="16"
          :width="3"
          :rotate="270"
          color="primary"
          :value="item.progress"
        />
        <span>{{ progressToPercentage(item.progress) }}</span>
      </template>
    </v-data-table>
    <v-list
      single-line
      class="width100"
    >
      <template v-if="data.markov.id !== null">
        <v-divider />
        <v-list-item class="px-2 py-1">
          <v-list-item-action class="pr-3 key">
            Markov:
          </v-list-item-action>
          <v-list-item-content>
            <v-list-item-title class="text-right">
              <a
                target="_blank"
                :href="$serverAddr + '/markovChains/' + data.markov.id"
              >{{
                data.markov.name }}</a>
            </v-list-item-title>
          </v-list-item-content>
        </v-list-item>
      </template>
    </v-list>
  </fc-tile>
</template>

<script>
  import tile from '@/components/tile/fc_tile.vue'
  export default {
    name: "MaskDetail",
    components: {
      'fc-tile': tile
    },
    props: {
      data: Object,
      customCharsets: {
        type: [Array, null],
        default: null
      }
    },
    data: function () {
      return {
        maskHeaders: [
          {
            text: 'Mask',
            align: 'left',
            value: 'mask'
          },
          {text: 'Progress', value: 'progress', align: 'right'}
        ],
        represenArray: {
          'l': {
            'chars': 'abcdefghijklmnopqrstuvwxyz',
            'represent': 'a-z'
          },
          'u': {
            'chars': 'ABCDEFGHIJKLMNOPQRSTUVWXYZ',
            'represent': 'A-Z'
          },
          'd': {
            'chars': '0123456789',
            'represent': '0-9'
          },
          'h': {
            'chars': '0123456789abcdef',
            'represent': '0-f'
          },
          'H': {
            'chars': '0123456789ABCDEF',
            'represent': '0-F'
          },
          's': {
            'chars': '«space»!"#$%&()*+,-./:;<=>?@[]^_`{|}~',
            'represent': 'special'
          },
          'a': {
            'chars': '?l?u?d?s',
            'represent': 'a-z,A-Z,0-9,special'
          },
          'b': {
            'chars': '0x00 - 0xff',
            'represent': 'ASCII'
          }
        }
      }
    },
    methods: {
      progressToPercentage: function (progress) {
        return progress.toFixed() + '%'
      },
      parseHTML(mask) {
        console.log(mask)
        var nextCharSymbol = false
        var parsed = ''
        for (var i = 0; i < mask.length; i++) {
          if (nextCharSymbol) {
            if (mask[i] === '?') {
              parsed += '?'
              nextCharSymbol = false
              continue
            }
            var symbol = this.represenArray[mask[i]]
            if (symbol === undefined) {
              symbol = {
                represent: mask[i]
              }
            }
            parsed += '<span class="symbol primary--text">' + symbol.represent + '</span>'
            nextCharSymbol = false
            continue
          }
          if (mask[i] === '?') {
            nextCharSymbol = true
          } else {
            parsed += mask[i]
            nextCharSymbol = false
          }
        }
        return parsed
      }
    }
  }
</script>

<style scoped>

</style>
<style >
  .symbol {
    margin-left: 1px;
    margin-right: 1px;
    border: 1px solid currentColor;
    padding: 3px;
    border-radius: 4px;
    font-weight: 500;
  }
</style>
