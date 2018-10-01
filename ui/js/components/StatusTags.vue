<template>
  <div class="field is-grouped is-grouped-multiline" style="padding-top:5px; padding-bottom:5px" v-if="thresholds">
    <div class="control" v-for="(value, key) in thresholds" v-if="value.threshold > 0" :key="key">
        <b-taglist attached>
          <b-tag type="is-dark">{{ key }} </b-tag>
          <b-tag :type="value.status === 'ok' ? 'is-success' : 'is-danger'">{{ formatFloat(value.threshold / 1000000) }} ms</b-tag>
        </b-taglist>
    </div>
  </div>
</template>

<script>
import common from './common.js'

export default {
  props: {
    thresholds: Object
  },
  mixins: [common],
}
</script>
