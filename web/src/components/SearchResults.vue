<template>
  <div class="results">
    <ul>
      <transition-group tag="li" name="results-list">
        <li v-for="result in results" :key="result">
          <p v-html="result.f"></p>
          <button @click.prevent="moreInfo(result)">More info...</button>
        </li>
      </transition-group>
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
  margin: 10px;
  padding: 10px;
  overflow-y: auto;
  overflow-x: hidden;
  height: 100%;
}

.results-list-enter-from {
  opacity: 0;
  transform: translateX(-30px);
}

.results-list-enter-active {
  transition: all 1s ease-out;
}

.results-list-enter-to {
  opacity: 1;
  transform: translateX(0);
}

.results-list-move {
  transition: transform 0.8s ease;
}
</style>