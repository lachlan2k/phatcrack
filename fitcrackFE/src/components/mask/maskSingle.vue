<!--
   * Author : see AUTHORS
   * Licence: MIT, see LICENSE
-->

<template>
  <div class="mask-container">
    <transition name="toolbar-fade">
      <v-sheet
        :id="id ? `${id}-chips-editor` : undefined"
        v-show="focused"
        class="charset-toolbar pa-1"
        :elevation="8"
      >
        <v-chip
          v-for="(item, key) in represenArray"
          :key="key"
          color="primary"
          class="ma-1"
          @click="addSymbol(key)"
        >
          {{ item.represent }}
        </v-chip>
        <v-chip
          v-for="(item, index) in customCharsets"
          :key="index"
          color="info"
          class="ma-1"
          @click="addSymbol(index + 1)"
        >
          {{ item.name }}
        </v-chip>
      </v-sheet>
    </transition>
    <div class="mask-input">
      <v-text-field
        :id="id ? `${id}-mask-input` : undefined"
        ref="field"
        label="Enter mask"
        filled
        outlined
        dense
        single-line
        :rules="maskRules"
        :value="value"
        @input="e => rawValue = e"
        @focus="focused = true"
        @blur="focused = false"
      />
      <v-btn
        :id="id ? `${id}-mask-remove` : undefined"
        v-show="!nonRemovable"
        class="ml-2 remove-button"
        icon
        @click="$emit('remove')"
      >
        <v-icon>mdi-delete</v-icon>
      </v-btn>
    </div>
  </div>
</template>

<script>
  export default {
    name: "MaskSingle",
    props: {
      id: {
        type: String,
        default: ''
      },
      value: {
        type: String,
        default: ''
      },
      nonRemovable: {
        type: Boolean,
        default: false
      },
      customCharsets: {
        type: [Array, null],
        default: null
      }
    },
    data() {
      return {
        rawValue: '',
        focused: false,
        maskRules: [
          v => /^(\?[ludhHsab\?]|[ -~])*$/.test(v) || 'Invalid mask',
          v => {
            let nextCharSymbol = false
            for (let c of v) {
              if (nextCharSymbol) {
                if (c === '?') {
                  nextCharSymbol = false
                  continue
                }
                let symbol = this.represenArray[c]
                if (symbol === undefined) {
                  if (this.customCharsets == null || !this.customCharsets[parseInt(c) - 1]) {
                    return `Unknown symbol ?${c}`
                  }
                }
                nextCharSymbol = false
                continue
              }
              if (c === '?') {
                nextCharSymbol = true
              } else {
                nextCharSymbol = false
              }
            }
            return true
          }
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
    computed: {
      /*
      parsedHTML() {
        let nextCharSymbol = false
        let parsed = ''
        for (let c of this.rawValue) {
          if (nextCharSymbol) {
            if (c === '?') {
              parsed += '?'
              nextCharSymbol = false
              continue
            }
            let symbol = this.represenArray[c]
            if (symbol === undefined) {
              if (this.customCharsets !== null && this.customCharsets[parseInt(c) - 1]) {
                symbol = {
                  represent: this.customCharsets[parseInt(c) - 1].name
                }
              } else {
                return '<span class="grey--text">Not valid mask.</span>'
              }
            }
            parsed += '<span class="symbol primary--text">' + symbol.represent + '</span>'
            nextCharSymbol = false
            continue
          }
          if (c === '?') {
            nextCharSymbol = true
          } else {
            parsed += c
            nextCharSymbol = false
          }
        }
        return parsed
      }
      */
    },
    watch: {
      rawValue (val) {
        this.$emit('input', val)
      }
    },
    methods: {
      addSymbol: function (symbol) {
        this.rawValue += '?' + symbol
        this.$refs.field.$el.value = this.rawValue
        this.$refs.field.focus()
      }
    },
  }
</script>
<style scoped>
.mask-container {
  position: relative;
}

.charset-toolbar {
  position: absolute;
  z-index: 10;
  bottom: 115%;
  max-width: 100%;
  user-select: none;
  border-radius: 1.3em;
  padding: 0;
}

.toolbar-fade-enter, .toolbar-fade-leave-to {
  opacity: 0;
  transform: translateY(1em);
}
.toolbar-fade-enter-active {
  transition: .35s cubic-bezier(0.19, 1, 0.22, 1);
}
.toolbar-fade-leave-active {
  transition: .15s .1s ease-in;
}

.mask-input {
  display: flex;
}
</style>
