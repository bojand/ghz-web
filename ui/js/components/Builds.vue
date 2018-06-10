<template>
  <section>
    <b-table 
      :data="data" 
      :loading="loading"

      paginated
      backend-pagination
      :total="total"
      :per-page="perPage"
      @page-change="onPageChange">

      <template slot-scope="props">
        <b-table-column field="id" label="ID" width="100">
          {{ props.row.id }}
        </b-table-column>

        <b-table-column field="name" label="Name" width="200">
          {{ props.row.name }}
        </b-table-column>

        <b-table-column field="description" label="Description">
          {{ props.row.description | truncate(80) }}
        </b-table-column>

        <b-table-column width="100">
          <button class="button block" @click="detailsClicked(props.row.id, $event)">Details</button>
        </b-table-column>
      </template>
    </b-table>
  </section>
</template>

<script>
import axios from 'axios'

export default {
  data() {
    return {
      data: [],
      total: 100,
      loading: false,
      page: 1,
      perPage: 20
    }
  },
  methods: {
    loadAsyncData() {
      const page = this.page - 1 || 0
      const params = `page=${page}`

      this.loading = true
      axios
        .get(`http://localhost:3000/api/projects?${params}`)
        .then(({ data }) => {
          console.log(data)
          this.data = data
          this.loading = false
        })
        .catch(error => {
          this.data = []
          this.total = 0
          this.loading = false
          throw error
        })
    },

    onPageChange(page) {
      this.page = page
      this.loadAsyncData()
    },

    onSort(field, order) {
      this.sortField = field
      this.sortOrder = order
      this.loadAsyncData()
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
  mounted() {
    this.loadAsyncData()
  }
}
</script>
