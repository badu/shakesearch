<template>
  <div class="results">
    <ul>
      <li v-for="result in results" :key="result" @click="moreInfo(result)">
        <p v-html="result.f"></p>
      </li>
    </ul>
  </div>
</template>
<script>
export default {
  props: ["results"],
  emits: ["more-info"],
  methods: {
    moreInfo(result) {
      this.$emit("more-info", result);
    },
  },
};
</script>
<style scoped>
.results {
  display: flex 1;
  background: #ffffff;
  opacity: 0.5;
  box-shadow: 0px 2px 6px 4px rgba(189, 170, 165, 0.1);
  border-radius: 9px;
  overflow: scroll;
}
</style>