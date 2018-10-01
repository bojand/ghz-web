<template>
  <nav class="navbar">
    
    <div class="navbar-brand">
      <a class="navbar-item" href="/">
        <!-- <img src="../../img/green_fwd_64.png" alt="Logo"> -->
        <span class="title is-3" style="padding-left: 12px; color: #3CBCBC;">
          <strong>ghz</strong>
        </span>
      </a>
      <span class="navbar-burger burger" data-target="navMenu">
        <span></span>
        <span></span>
        <span></span>
      </span>
    </div>
    
    <div id="navMenu" class="navbar-menu">
      <div class="navbar-end">
        <div class="navbar-item">
          <div class="field is-grouped">
            <p class="control">
              <button class="button" @click="infoModal()">
                Info
              </button>
            </p>
            <p class="control">
              <a class="button block" href="https://bulma.io/">
                Docs
              </a>
            </p>
            <p class="control">
              <a class="button block" href="https://bulma.io/">
                API
              </a>
            </p>
          </div>
        </div>
      </div>
    </div>

  </nav>
</template>

<script>
import InfoModal from '../components/InfoModal.vue'

export default {
  methods: {
    async infoModal() {
      let info = await this.$store.fetchInfo()

      if (typeof info === 'string') {
        info = JSON.parse(info)
      }

      info = JSON.stringify(info, null, 2)

      console.log(info)

      this.$modal.open({
        parent: this,
        component: InfoModal,
        hasModalCard: true,
        props: {
          info
        }
      })
    }
  }
}
</script>
