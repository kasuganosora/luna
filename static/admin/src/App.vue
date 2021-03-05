<template>
  <div id="app">
    <adminPageLayout v-if="layoutType === 'full'"/>
    <emptyLayout v-if="layoutType === 'empty'"/>
  </div>
</template>

<script>
import adminPageLayout from "@/views/layout";
import emptyLayout from "@/views/empty-layout";
export default {
  name: 'app',
  components: {
    adminPageLayout,
    emptyLayout,
  },
  computed: {
    layoutType(){
      if(this.$route.meta.layout !== undefined){
        return this.$route.meta.layout
      }
      return "full"
    },
  }
}
</script>

