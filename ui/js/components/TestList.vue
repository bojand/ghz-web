<template>
  <section>
    <h2 class="subtitle"><strong>Tests</strong></h2>
    <b-table 
      :data="data" 
      :loading="loading"

      paginated
      backend-pagination
      :total="total"
      :per-page="perPage"
      @page-change="onPageChange"

      :opened-detailed="defaultOpenedDetails"
      detailed
      detail-key="id"

      backend-sorting
      :default-sort-direction="defaultSortOrder"
      :default-sort="[sortField, sortOrder]"
      @sort="onSort">

      <template slot-scope="props">
        <b-table-column field="id" label="ID" width="100" sortable>
          {{ props.row.id }}
        </b-table-column>

        <b-table-column field="name" label="Name" width="200" sortable>
          {{ props.row.name }}
        </b-table-column>

        <b-table-column field="description" label="Description">
          {{ props.row.description | truncate(80) }}
        </b-table-column>

        <b-table-column field="status" label="Status" centered>
          <b-icon 
            :icon="props.row.status === 'ok' ? 'checkbox-marked-circle-outline' : 'alert-circle-outline'" 
            size="is-medium"
            custom-size="mdi-24px"
            :type="props.row.status === 'ok' ? 'is-success' : 'is-danger'"
          ></b-icon>
        </b-table-column>

        <b-table-column width="100">
          <router-link :to="{ name: 'test', params: { projectId: projectId, testId: props.row.id } }" class="button is-info">Details</router-link>
        </b-table-column>
      </template>

      <template slot="detail" slot-scope="props">
        <component-test-details :project-id="projectId" :testId="props.row.id"></component-test-details>
      </template>
    </b-table>
  </section>
</template>

<script>
import axios from 'axios'
import TestRowDetails from './TestRowDetails.vue'

export default {
  data() {
    return {
      data: [],
      total: 100,
      loading: false,
      sortField: 'id',
      sortOrder: 'desc',
      defaultSortOrder: 'desc',
      page: 1,
      perPage: 20,
      defaultOpenedDetails: []
    }
  },
  props: {
    projectId: [String, Number]
  },
  watch: {
    projectId(newVal, oldVal) {
      this.loadData()
    }
  },
  mounted() {
    this.loadData()
  },
  methods: {
    async loadData() {
      const page = this.page - 1 || 0
      const params = `page=${page}&sort=${this.sortField}&order=${this.sortOrder}`

      this.loading = true
      try {
        const { data } = await axios.get(
          `http://localhost:3000/api/projects/${this.projectId}/tests?${params}`
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
    'component-test-details': TestRowDetails
  }
}
</script>
