<template>
  <section>
    <div class="modal-card" style="width: auto">
      <header class="modal-card-head">
        <p class="modal-card-title">About</p>
      </header>
      <section class="modal-card-body">
        <p>
          <b-message>
            <pre style="background-color: transparent; white-space: pre-wrap;">{{ info }}</pre>
          </b-message>
        </p>
      </section>
      <footer class="modal-card-foot">
        <button class="button" type="button" @click="$parent.close()">Close</button>
      </footer>
    </div>
  </section>
</template>

<script>
export default {
  props: {
    info: String
  },
  mounted() {
    console.log('mounted')
    console.log(this.info)
  }
}
</script>
