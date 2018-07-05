<template>
  <section>
    <h2 class="subtitle strong"><strong>Test Details</strong></h2>
    <div class="box">

      <div class="content" v-if="!editMode">
        
        <span class="title is-4">
          <strong>{{ model.name }}</strong>
          <b-icon 
            :icon="model.status === 'ok' ? 'checkbox-marked-circle-outline' : 'alert-circle-outline'" 
            size="is-medium"
            custom-size="mdi-18px"
            :type="model.status === 'ok' ? 'is-success' : 'is-danger'"
          ></b-icon>
        </span>
        <div class="level-left" style="padding-top:5px; padding-bottom:10px" v-if="model.thresholds">
          <div class="level-item" v-for="(value, key) in model.thresholds" :key="key">
            <div class="control">
              <b-taglist attached>
                  <b-tag type="is-info">{{ key }} </b-tag>
                  <b-tag :type="value.status === 'ok' ? 'is-success' : 'is-danger'">{{ value.threshold }} ms</b-tag>
              </b-taglist>
            </div>
            <!-- {{ key }}: {{ value.threshold }}
            <b-icon 
              :icon="value.status === 'ok' ? 'checkbox-marked-circle-outline' : 'alert-circle-outline'" 
              size="is-small"
              custom-size="mdi-18px"
              :type="value.status === 'ok' ? 'is-success' : 'is-danger'"
            ></b-icon> -->
          </div>
        </div>
        
        <div>
          <p>
            {{ model.description }}
          </p>
        </div>

      </div>

      <div class="content" v-if="editMode">
        <div class="media">
          <div class="media-left">
            <b-field>
              <b-input placeholder="name" v-model="model.name" required></b-input>
            </b-field>
            <b-field>
              <b-input placeholder="description" v-model="model.description"></b-input>
            </b-field>
            <b-field>
              <b-select placeholder="Threshold" v-model="selectedThreshold" @input="selectedChaged">
                <option>median</option>
                <option>mean</option>
                <option>95th</option>
                <option>99th</option>
              </b-select>
              <b-input type="number" min="0" placeholder="0" v-model="selectedThresholdValue"></b-input>
              <p class="control">
                <button class="button is-success" @click="setThreshold">Set</button>
              </p>
          </b-field>
          <div class="field">
              <b-tooltip label="Fail the test if last run has any errors" position="is-right" type="is-light">
                <b-switch v-model="model.failOnError" type="is-danger">Fail on error</b-switch>
              </b-tooltip>
            </div>
          </div>
        </div>
      </div>

      <nav class="level">
        <div class="level-left">
          <a class="level-item" aria-label="reply">
            <button :class="['button', editMode ? 'is-primary' : '']" @click="editClicked">
              <b-icon :icon="editMode ? 'check' : 'pencil'" size="is-small"></b-icon>
              <span>{{ editMode ? 'Save' : 'Edit' }}</span>
            </button>
          </a>
          <a v-if="editMode" class="level-item" aria-label="reply">
            <button class="button" @click="cancelClicked">
              <b-icon icon="cancel" size="is-small"></b-icon>
              <span>Cancel</span>
            </button>
          </a>
        </div>
      </nav>

    </div>
  </section>
</template>

<script>
import axios from 'axios'

export default {
  data() {
    return {
      id: null,
      loading: false,
      editMode: false,
      model: {
        name: '',
        description: '',
        failOnError: false,
        thresholds: null
      },
      selectedThreshold: 'median',
      selectedThresholdValue: 0
    }
  },
  props: {
    projectId: [String, Number],
    testId: [String, Number]
  },
  watch: {
    projectId(newVal, oldVal) {
      this.loadData()
    },
    testId(newVal, oldVal) {
      this.loadData()
    }
  },
  store: ['project', 'test'],
  methods: {
    async loadData() {
      this.loading = true
      try {
        if (!this.test) {
          this.test = await this.$store.fetchTest(this.projectId, this.testId)
        }

        Object.assign(this.model, this.test)
        if (!this.model.thresholds) {
          this.model.thresholds = {
            median: { status: 'ok', threshold: 10000 },
            mean: { status: 'fail', threshold: 20000 }
          }
        }
        this.loading = false
      } catch (e) {
        this.loading = false

        this.$snackbar.open({
          message: e.message,
          type: 'is-danger',
          position: 'is-top'
        })
      }
    },

    async editClicked() {
      if (this.editMode) {
        this.loading = true

        try {
          this.test = await this.$store.updateTest(this.projectId, this.model)
          Object.assign(this.model, this.test)

          this.loading = false
        } catch (e) {
          this.loading = false

          this.$snackbar.open({
            message: e.message,
            type: 'is-danger',
            position: 'is-top'
          })

          this.loadData()
        }
      }
      this.editMode = !this.editMode
    },

    async cancelClicked() {
      await this.loadData()
      this.editMode = false
    },

    setThreshold() {
      if (!this.model.thresholds) {
        this.model.thresholds = {
          mean: { status: 'ok', threshold: 0 },
          median: { status: 'ok', threshold: 0 },
          '96th': { status: 'ok', threshold: 0 },
          '99th': { status: 'ok', threshold: 0 }
        }
      }

      if (!this.model.thresholds[this.selectedThreshold]) {
        this.model.thresholds[this.selectedThreshold] = { status: 'ok', threshold: 0 }
      }

      const val = parseInt(this.selectedThresholdValue, 10)
      this.model.thresholds[this.selectedThreshold].threshold = val
    },

    selectedChaged() {
      if (!this.model.thresholds) {
        this.selectedThresholdValue = 0
        return
      }

      if (!this.model.thresholds[this.selectedThreshold]) {
        this.selectedThresholdValue = 0
        return
      }

      this.selectedThresholdValue = this.model.thresholds[this.selectedThreshold].threshold || 0
    }
  },
  mounted() {
    this.loadData()

    this.selectedChaged()
  }
}
</script>
