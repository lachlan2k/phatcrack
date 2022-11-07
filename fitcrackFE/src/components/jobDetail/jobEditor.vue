<template>
  <div>
    <!-- name -->
    <v-list-item>
      <v-list-item-content>
        <v-text-field
          v-model="edit.name"
          outlined
          hide-details
        >
          <template #label>
            <v-icon class="label-icon">
              mdi-format-title
            </v-icon>
            Name
          </template>
        </v-text-field>
      </v-list-item-content>
    </v-list-item>
    <!-- job start -->
    <v-list-item>
      <v-list-item-content class="py-0">
        <v-list-item-title class="d-flex align-center">
          <div class="subtitle-2">
            <v-icon>
              mdi-ray-start
            </v-icon>
            Start
          </div>
          <v-spacer />
          <v-checkbox
            v-model="edit.startNow"
            label="No Start Delay"
          />
        </v-list-item-title>
        <picker
          v-show="!edit.startNow"
          v-model="edit.time_start"
          outlined
          hide-details
        />
      </v-list-item-content>
    </v-list-item>
    <!-- job end -->
    <v-list-item>
      <v-list-item-content class="pt-0">
        <v-list-item-title class="d-flex align-center">
          <div class="subtitle-2">
            <v-icon>
              mdi-ray-end
            </v-icon>
            End
          </div>
          <v-spacer />
          <v-checkbox
            v-model="edit.endNever"
            label="No Time Limit"
          />
        </v-list-item-title>
        <picker
          v-show="!edit.endNever"
          v-model="edit.time_end"
          outlined
          hide-details
        />
      </v-list-item-content>
    </v-list-item>
    <!-- wu secs -->
    <v-list-item>
      <v-list-item-content>
        <v-text-field
          v-model="edit.seconds_per_job"
          outlined
          hide-details
          type="number"
          suffix="seconds"
        >
          <template #label>
            <v-icon class="label-icon">
              mdi-timeline-clock
            </v-icon>
            Time per workunit
          </template>
        </v-text-field>
      </v-list-item-content>
    </v-list-item>
    <!-- comment -->
    <v-list-item>
      <v-list-item-content>
        <v-textarea
          v-model="edit.comment"
          outlined
          auto-grow
          hide-details
        >
          <template #label>
            <v-icon class="label-icon">
              mdi-format-quote-close
            </v-icon>
            Comment
          </template>
        </v-textarea>
      </v-list-item-content>
    </v-list-item>
    <!-- prince copypasta -->
    <template v-if="data.attack_mode == 8">
      <v-list-item>
        <v-list-item-content>
          <v-checkbox
            v-model="edit.check_duplicates"
            label="Check password duplicates"
          />
          <v-checkbox
            v-model="edit.case_permute"
            label="Case permutation"
          />
        </v-list-item-content>
      </v-list-item>

      <v-list-item>
        <v-list-item-content>
          <v-text-field
            v-model="edit.min_password_len"
            outlined
            label="Minimal length of passwords (1 - 32)"
            required
            type="number"
            min="1"
            max="32"
          />
        </v-list-item-content>
      </v-list-item>

      <v-list-item>
        <v-list-item-content>
          <v-text-field
            v-model="edit.max_password_len"
            outlined
            label="Maximal length of passwords (1 - 32)"
            required
            type="number"
            min="1"
            max="32"
          />
        </v-list-item-content>
      </v-list-item>

      <v-list-item>
        <v-list-item-content>
          <v-text-field
            v-model="edit.min_elem_in_chain"
            outlined
            label="Minimal number of elements per chain (1 - 16)"
            required
            type="number"
            min="1"
            max="16"
          />
        </v-list-item-content>
      </v-list-item>

      <v-list-item>
        <v-list-item-content>
          <v-text-field
            v-model="edit.max_elem_in_chain"
            outlined
            label="Maximal number of elements per chain (1 - 16)"
            required
            type="number"
            min="1"
            max="16"
          />
        </v-list-item-content>
      </v-list-item>

      <v-list-item>
        <v-list-item-content>
          <v-text-field
            v-model="edit.generate_random_rules"
            outlined
            label="Generate random rules"
            required
            type="number"
            min="0"
          />
        </v-list-item-content>
      </v-list-item>
    </template>
    <!-- save button -->
    <v-list-item>
      <v-list-item-content>
        <v-btn
          color="primary"
          @click="save"
        >
          <v-icon left>
            mdi-content-save
          </v-icon>
          <span>Save</span>
        </v-btn>
      </v-list-item-content>
    </v-list-item>
  </div>
</template>

<script>
import Picker from '@/components/picker/datetime.vue'

export default {
  components: {
    Picker
  },
  props: {
    data: {
      type: Object,
      default: () => {}
    }
  },
  data () {
    return {
      edit: {
        name: this.data.name,
        comment: this.data.comment,
        seconds_per_job: this.data.seconds_per_job,
        time_start: this.data.time_start === null ?
        this.$moment().format('YYYY-MM-DDTHH:mm') :
        this.$moment.utc(this.data.time_start).local().format('YYYY-MM-DDTHH:mm'),
        time_end: this.data.time_end === null ?
          this.$moment().format('YYYY-MM-DDTHH:mm') :
          this.$moment.utc(this.data.time_end).local().format('YYYY-MM-DDTHH:mm'),
        startNow: (this.data.time_start === null),
        endNever: (this.data.time_end === null),
        check_duplicates: this.data.check_duplicates,
        case_permute: this.data.case_permute,
        min_password_len: this.data.min_password_len,
        max_password_len: this.data.max_password_len,
        min_elem_in_chain: this.data.min_elem_in_chain,
        max_elem_in_chain: this.data.max_elem_in_chain,
        generate_random_rules: this.data.generate_random_rules
      }
    }
  },
  methods: {
    save () {
      
      if (this.edit.startNow) {
        this.edit.time_start = ''
      } else {
        this.edit.time_start = this.$moment(this.edit.time_start).utc().toISOString(true).slice(0, 16)
      }

      if (this.edit.endNever) {
        this.edit.time_end = ''
      } else {
        this.edit.time_end = this.$moment(this.edit.time_end).utc().toISOString(true).slice(0, 16)
      }

      this.axios.put(this.$serverAddr + '/job/' + this.data.id , this.edit
      ).then((response) => {
        this.$emit('save')
      })
    }
  }
}
</script>

<style scoped>
.label-icon {
  color:inherit;
  transition:none;
}
</style>