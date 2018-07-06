<template>
  <div class="level-left" style="padding-top:5px; padding-bottom:5px" v-if="thresholds">
    <div class="level-item" v-for="(value, key) in thresholds" :key="key">
      <div class="control">
        <b-taglist attached>
          <b-tag type="is-dark">{{ key }} </b-tag>
          <b-tag :type="value.status === 'ok' ? 'is-success' : 'is-danger'">{{ value.threshold }} ms</b-tag>
        </b-taglist>
      </div>
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
