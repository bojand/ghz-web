<template>
  <section>
    <h2 class="subtitle"><strong>Runs</strong></h2>
     <b-table 
      :data="data" 
      :loading="loading"

      paginated
      backend-pagination
      :total="total"
      :per-page="perPage"
      @page-change="onPageChange"

      backend-sorting
      :default-sort-direction="defaultSortOrder"
      :default-sort="[sortField, sortOrder]"
      @sort="onSort">

      <template slot-scope="props">
        <!-- <b-table-column field="id" label="ID" sortable>
          {{ props.row.id }}
        </b-table-column> -->

        <b-table-column field="date" label="Date" sortable>
          <router-link :to="{ name: 'run', params: { projectId: projectId, testId: testId, runId: props.row.id } }">
            {{ new Date(props.row.date).toLocaleString() }}
            <!-- <b-icon icon="open-in-app" size="is-small"></b-icon> -->
          </router-link>
        </b-table-column>

        <b-table-column field="count" label="Count" sortable>
          {{ props.row.count }}
        </b-table-column>

        <b-table-column field="total" label="Total" sortable>
          {{ formatNano(props.row.total) }} ms
        </b-table-column>

        <b-table-column field="average" label="Average" sortable>
          <span class="tag" :class="classifyResult(test, props.row.average, 'mean')">
            {{ formatNano(props.row.average) }} ms
          </span>
        </b-table-column>

        <b-table-column field="slowest" label="Slowest" sortable>
          <span class="tag" :class="classifyResult(test, props.row.slowest, 'slowest')">
            {{ formatNano(props.row.slowest) }} ms
          </span>
        </b-table-column>

        <b-table-column field="fastest" label="Fastest" sortable>
          <span class="tag" :class="classifyResult(test, props.row.fastest, 'fastest')">
            {{ formatNano(props.row.fastest) }} ms
          </span> 
        </b-table-column>

        <b-table-column field="rps" label="RPS" sortable>
          <span class="tag" :class="classifyResult(test, props.row.rps, 'RPS')">
            {{ formatFloat(props.row.rps, 0) }}
          </span> 
        </b-table-column>

        <b-table-column field="status" label="Status" centered>
          <b-icon 
            :icon="props.row.status === 'ok' ? 'checkbox-marked-circle-outline' : 'alert-circle-outline'" 
            size="is-medium"
            custom-size="mdi-24px"
            :type="props.row.status === 'ok' ? 'is-success' : 'is-danger'"
          ></b-icon>
        </b-table-column>

        <!-- <b-table-column>
          <router-link :to="{ name: 'run', params: { projectId: projectId, testId: testId, runId: props.row.id } }" class="button is-info">Details</router-link>
        </b-table-column> -->
      </template>

      <!-- <template slot="detail" slot-scope="props">
        <component-run-details :project-id="projectId" :testId="testId" :runId="props.row.id"></component-run-details>
      </template> -->
    </b-table>
  </section>
</template>

<script>
import axios from 'axios'
import RunRowDetails from './RunRowDetails.vue'
import common from './common.js'

export default {
  data() {
    return {
      data: [],
      total: 100,
      loading: false,
      sortField: 'date',
      sortOrder: 'desc',
      defaultSortOrder: 'desc',
      page: 1,
      perPage: 20,
      defaultOpenedDetails: []
    }
  },
  props: {
    projectId: [String, Number],
    testId: [String, Number]
  },
  store: ['test'],
  watch: {
    projectId(newVal, oldVal) {
      this.loadData()
    },
    testId(newVal, oldVal) {
      this.loadData()
    }
  },
  mounted() {
    this.loadData()
  },
  mixins: [common],
  methods: {
    async loadData() {
      const page = this.page - 1 || 0
      const params = `page=${page}&sort=${this.sortField}&order=${this.sortOrder}`

      this.loading = true
      try {
        const { data } = await axios.get(
          `http://localhost:3000/api/projects/${this.projectId}/tests/${this.testId}/runs?${params}`
        )

        this.data = data.data
        this.total = data.total
        this.loading = false
      } catch (e) {
        this.data = []
        this.total = 0
        this.loading = false

        this.$snackbar.open({
          message: e.message,
          type: 'is-danger',
          position: 'is-top'
        })
      }
    },

    onPageChange(page) {
      this.page = page
      this.loadData()
    },

    onSort(field, order) {
      this.sortField = field
      this.sortOrder = order
      this.loadData()
    },

    detailsClicked(id, ev) {
      console.log(id)
    }
  },
  filters: {
    /**
     * Filter to truncate string, accepts a length parameter
     */
    truncate(value, length) {
      return value.length > length ? value.substr(0, length) + '...' : value
    }
  },
  components: {
    'component-run-details': RunRowDetails
  }
}
</script>
