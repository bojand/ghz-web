<template>
  <div class="field is-grouped is-grouped-multiline" style="padding-top:5px; padding-bottom:5px" v-if="thresholds">
    <div class="control" v-for="(value, key) in thresholds" v-if="value.threshold > 0" :key="key">
        <b-taglist attached>
          <b-tag type="is-dark">{{ key }} </b-tag>
          <b-tag :type="value.status === 'ok' ? 'is-success' : 'is-danger'">{{ value.threshold }} ms</b-tag>
        </b-taglist>
    </div>
  </div>
</template>

<script>
export default {
  props: {
    thresholds: Object
  }
}
</script>
