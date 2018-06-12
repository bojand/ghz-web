<template>
  <section>
    <component-project-details></component-project-details>
    <hr />
    <section>
      <h2 class="subtitle"><strong>Tests</strong></h2>
      <b-field grouped>
      </b-field>
      <b-table 
        :data="data" 
        :loading="loading"

        :default-sort-direction="defaultSortDirection"
        :default-sort="[sortField, sortOrder]"

        paginated
        backend-pagination
        :total="total"
        :per-page="perPage"
        @page-change="onPageChange">

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

          <b-table-column width="100">
            <router-link :to="{ name: 'project', params: { id: props.row.id } }" class="button block">Details</router-link>
          </b-table-column>
        </template>
      </b-table>
    </section>
  </section>
</template>

<script>
import ProjectDetails from './ProjectDetails.vue'

export default {
  data() {
    return {
      data: [],
      total: 100,
      loading: false,
      sortField: 'id',
      sortOrder: 'desc',
      defaultSortDirection: 'desc',
      page: 1,
      perPage: 20
    }
  },
  methods: {
    onPageChange(page) {
      this.page = page
    },

    onSort(field, order) {
      this.sortField = field
      this.sortOrder = order
    },
  },
  components: {
    'component-project-details': ProjectDetails
  }
}
</script>
