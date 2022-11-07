<!--
   * Author : see AUTHORS
   * Licence: MIT, see LICENSE
-->

<template>
  <div class="textAreaCont">
    <v-row 
      v-if="searchEnable"
      class="mx-2 mb-2"
    >
      <v-text-field
        v-model="searchText"
        single-line
        hide-details
        label="Search"
        class="ml-3"
        clearable
        @click:clear="clearSearch"
        @keypress.enter="search"
      />
      <v-btn
        icon
        color="primary"
        class="mt-2"
        @click="search"
      >
        <v-icon>search</v-icon>
      </v-btn>
    </v-row>
    <div
      class="scrollCont"
      :style="style"
      :class="{editable: !readonly}"
    >
      <div
        v-if="noResults"
        class="text-center pt-4"
      >
        No result found.
      </div>
      <div
        v-else-if="loading"
        class="text-center pa-5"
      >
        <v-progress-circular
          slot="spinner"
          size="50"
          :width="3"
          indeterminate
          color="primary"
        />
      </div>
      <v-row
        v-else
        class="minheight200 mx-2"
      >
        <div class="lineNumberCont pr-2 pl-3 pt-2">
          <p
            v-for="index in linesCount"
            :key="index"
            class="pa-0 ma-0"
          >
            {{ index }}
          </p>
        </div>
        <div class="textarea-wrap">
          <textarea
            v-if="searching"
            v-model="searchData"
            class="pt-2 pl-2"
            :readonly="true"
            wrap="off"
          />
          <textarea
            v-else
            ref="textarea"
            v-model="textareaData"
            class="pt-2 pl-2"
            :readonly="readonly"
            wrap="off"
            @input="update"
            @blur="emmitBlur"
            @focus="emmitFocus"
          />
          <infinite-loading
            v-if="url && !searching"
            @infinite="loadData"
          >
            <v-progress-circular
              slot="spinner"
              size="50"
              :width="3"
              indeterminate
              color="primary"
            />
            <span slot="no-more">
              You reached the end of file.
            </span>
            <span slot="no-results">
              You reached the end of file.
            </span>
          </infinite-loading>
        </div>
        <slot name="after" />
        <div
          v-if="canRemoveLine"
          class="pt-2"
        >
          <p
            v-for="index in linesCount"
            :key="index"
            class="pa-0 ma-0"
          >
            <v-icon
              small
              color="error"
              class="removeLineBtn"
              @click="removeLine(index)"
            >
              close
            </v-icon>
          </p>
        </div>
      </v-row>
    </div>
  </div>
</template>

<script>
  import InfiniteLoading from 'vue-infinite-loading';
  export default {
    name: "FcTextarea",
    components: {
      InfiniteLoading,
    },
    props: {
      'maxHeight': String,
      'readonly': Boolean,
      'url': String,
      'searchEnable': Boolean,
      'canRemoveLine': Boolean,
      'value': {
        type: String,
        default: ''
      },
    },
    data() {
      return {
        data: '',
        page: 1,
        searching: false,
        searchText: '',
        searchData: '',
        noResults: false,
        loading: false,
        textareaData: ''
      }
    },
    computed: {
      linesCount: function () {
        if (this.searching) {
          return this.searchData.split('\n').length - 1
        } else {
          if (this.value) {
            return this.value.split('\n').length
          } else {
            return 1
          }
        }
      },
      style () {
        if (this.maxHeight) {
          return 'max-height: ' + this.maxHeight + 'px'
        } else {
          return ''
        }
      }
    },
    watch:{
      value: function(){
        this.textareaData = this.value
      }
    },
    mounted: function () {
      this.textareaData = this.value
    },
    methods: {
      focus () {
        this.$refs.textarea.focus()
      },
      emmitBlur() {
        this.$emit('blur', this.value)
      },
      emmitFocus() {
        this.$emit('focus', this.value)
      },
      update() {
        if (!this.readonly) {
          this.$emit('input', this.textareaData)
        }
      },
      loadData: function ($state) {
        this.axios.get(this.url,
          {
            params: {
              'page': this.page,
              'per_page': 100
            }
          }).then((response) => {

          if (response.data.status) {
            this.value += response.data.data
            this.page++
            $state.loaded()
          }
          else {
            $state.complete()
          }
        });
      },
      search: function() {
        this.noResults = false
        this.loading = true
        this.searching = true
        this.searchData = ''
        this.axios.get(this.url,
          {
            params: {
              'search': this.searchText
            }
          }).then((response) => {
            if (!response.data.status) {
              this.noResults = true
              this.searchData = ''
            } else {
              this.searchData = response.data.data
            }
            this.loading = false
          });
      },
      clearSearch: function () {
        this.searching = false
        this.noResults = false
      },
      removeLine: function (index) {
        var lines = this.textareaData.split('\n')
        lines.splice(index - 1, 1)
        this.textareaData = lines.join('\n')
        this.$emit('input', this.textareaData)
        this.$emit('blur', this.textareaData)
      }
    },
  }
</script>

<style scoped>
  textarea{
    width: 100%;
    height: calc(100% - 20px);
    color: #222;
    resize: none;
    overflow-x: auto;
    overflow-y: hidden;
    background-image: linear-gradient(to bottom, rgba(0,0,0,.05) 50%, transparent 50%);
    background-size: 100% 3.199em;
    background-repeat: repeat-y;
    background-position-y: -1em;
  }
  .theme--dark textarea {
    color: #dedede;
    background-image: linear-gradient(to bottom, rgba(255,255,255,.05) 50%, transparent 50%);
  }

  .lineNumberCont {
    text-align: end;
    font-weight: bold;
    background-color: rgba(0,0,0,.1);
  }

  .textAreaCont {
    box-sizing: border-box;
  }
  .theme--dark .lineNumberCont { background-color: rgba(255,255,255,.1) }

  .scrollCont {
    height: 100%;
    overflow: auto;
    min-height: 200px;
    font-family: 'Roboto Mono', monospace;
    font-size: .85em;
    line-height: 1.6em;
  }

  .scrollCont.editable {
    border: 1px solid rgba(0,0,0,.15);
    border-radius: 4px;
  }
  .theme--dark .scrollCont { border-color: rgba(255,255,255,.15) }

  .textarea-wrap {
    flex: 1 1 auto;
  }

  .minheight200 {
    min-height: 200px;
  }

  .removeLineBtn {
    cursor: pointer;
  }

</style>
