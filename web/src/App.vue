<template>
  <div class="container-fluid overflow-hidden main-div">
    <!-- Left content -->
    <div class="search-box">
      <TheNavBar />
      <SearchBar :searchTerm="searchTerm" @search="performSearch" />
      <Shakespeare v-if="start" />
    </div>

    <!-- Right content -->
    <div class="search-result">
      <transition name="fade" mode="out-in">
        <SearchResults v-if="!start" :results="results" />
      </transition>
    </div>
  </div>
</template>

<script>
import TheNavBar from "./components/TheNavBar.vue";
import SearchBar from "./components/SearchBar.vue";
import Shakespeare from "./components/Shakespeare.vue";
import SearchResults from "./components/SearchResults.vue";

export default {
  components: {
    TheNavBar,
    SearchBar,
    Shakespeare,
    SearchResults,
  },
  computed: {
    start() {
      return this.$store.state.start;
    },
    searchTerm() {
      return this.$store.state.searchTerm;
    },
    results() {
      return this.$store.state.results;
    },
  },
  methods: {
    performSearch(searchTerm) {
      this.$store.dispatch("fetchResults", searchTerm);
    },
  },
};
</script>

<style>
.main-div {
  display: flex;
  backface-visibility: hidden;
  height: 100%;
  left: 0px;
  opacity: 1;
  overflow: scroll;
  position: absolute;
  top: 0px;
  width: 100%;
}
.search-box{
    flex:1;
    flex-direction: column;
    justify-content: center
}

.search-result{
    flex:1;
    flex-direction: column;
    justify-content: center
}

</style>
