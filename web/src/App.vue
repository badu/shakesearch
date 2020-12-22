<template>
  <div class="container-fluid overflow-hidden main-div">
    <!-- Left content -->
    <div class="search-box">
      <TheNavBar :lastError="lastError" />
      <SearchBar
        :searchTerm="searchTerm"
        @search="performSearch"
        @clearSearch="clearSearch"
      />
      <Shakespeare v-if="!results" />
    </div>

    <!-- Right content -->
    <div class="search-result">
      <transition name="fade" mode="out-in">
        <SearchResults v-if="results" :results="results" @moreInfo="moreInfo" />
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
    searchTerm() {
      return this.$store.state.searchTerm;
    },
    results() {
      return this.$store.state.results;
    },
    lastError() {
      return this.$store.state.lastError;
    },
  },
  methods: {
    performSearch(searchTerm) {
      this.$store.dispatch("fetchResults", searchTerm);
    },
    clearSearch() {
      if (this.$store.searchTerm === "") {
        return;
      }
      this.$store.dispatch("clearResults");
    },
    moreInfo(payload) {
      console.log("should more info", payload);
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
.search-box {
  flex: 1;
  flex-direction: column;
  justify-content: center;
}

.search-result {
  flex: 1;
  flex-direction: column;
  justify-content: center;
}

.animate {
  /* transform: translateX(-150px); */
  animation: slide-fade 0.3s ease-out forwards;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

.fade-enter-active {
  transition: opacity 0.3s ease-out;
}

.fade-leave-active {
  transition: opacity 0.3s ease-in;
}

.fade-enter-to,
.fade-leave-from {
  opacity: 1;
}

@keyframes slide-scale {
  0% {
    transform: translateX(0) scale(1);
  }

  70% {
    transform: translateX(-120px) scale(1.1);
  }

  100% {
    transform: translateX(-150px) scale(1);
  }
}
</style>
