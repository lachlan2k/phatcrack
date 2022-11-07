<!--
   * Author : see AUTHORS
   * Licence: MIT, see LICENSE
-->

<template>
  <div class="node">
    <span
      v-if="keyName"
      class="key"
    >
      <v-icon
        v-if="isObject"
        @click.stop="toggleNode"
      >{{ value ? 'arrow_drop_down' : 'arrow_drop_up' }}</v-icon>{{ keyName }}:
    </span>
    <commponent
      :is="`Json${valueType}`"
      v-model="toggle"
      :json-value="value"
      :key-name="keyName"
    />
  </div>
</template>

<script>
import JsonString from '@/components/json/types/json-string.vue';
import JsonNumber from '@/components/json/types/json-number.vue';
import JsonBoolean from '@/components/json/types/json-boolean.vue';
import JsonObject from '@/components/json/types/json-object.vue';
import JsonArray from '@/components/json/types/json-array.vue';

export default {
    name: 'JsonBox',
    components: {
        JsonString,
        JsonNumber,
        JsonBoolean,
        JsonObject,
        JsonArray
    },
    props: {
        value: [Object, Array, String, Number, Boolean],
        keyName: String
    },
    data() {
        return {
            toggle: true
        };
    },
    computed: {
        valueType() {
            if (Array.isArray(this.value)) {
                return 'Array';
            } else if (this.value === null){
              return 'String'
            } else if (typeof this.value === 'object') {
                return 'Object';
            } else if (typeof this.value === 'number') {
                return 'Number';
            } else if (typeof this.value === 'string') {
                return 'String';
            } else if (typeof this.value === 'boolean') {
                return 'Boolean';
            }
        },
        isObject() {
            return this.valueType == 'Array' || this.valueType == 'Object'; // eslint-disable-line
        }
    },
    methods: {
        toggleNode() {
            this.toggle = !this.toggle;
        }
    }
};
</script>

<style lang="scss">
.node {
    font-size: 14px;
    position: relative;
    color: #525252;
    font-family: Consolas,Menlo,Courier,monospace;
    white-space: nowrap;

    &:after {
        content: ','
    }
    &:last-of-type {
        &:after {
            content: ''
        }
    }

    & .node {
        margin-left: 25px;
    }

}
</style>
