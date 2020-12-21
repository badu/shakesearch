<template>
  <div class="position-relative vh-100">
    <div class="position-absolute top-45 start-49 translate-middle">
      <div class="row">
        <div class="col-1">
          <img
            src="../assets/svg/magnifying-glass.svg"
            class="py-1 px-2"
            alt="sketch"
            @click="performSearch"
          />
        </div>
        <div class="col">
          <input
            type="text"
            class="search-input"
            :value="searchTerm"
            v-on:keyup.enter="performSearch"
            placeholder="What art thee looking f'r?"
          />
        </div>
      </div>
      <img src="../assets/svg/signature.svg" alt="sketch" />
    </div>
  </div>
</template>

<script>
export default {
  props: ["searchTerm"],
  emits: ["search"],
  methods: {
    performSearch(event) {
      if (event.target.value === "") {
        event.target.placeholder = "Art thee looking f'r something?";
        return;
      }
      this.$emit("search", event.target.value);
    },
  },
};
</script>

<style scoped>
@import url("https://fonts.googleapis.com/css2?family=Courier+Prime:ital,wght@0,400;0,700;1,400;1,700&display=swap");
/*font-family: 'Courier Prime', monospace;*/

@import url("https://fonts.googleapis.com/css2?family=Playfair+Display:ital,wght@0,400;0,500;0,600;0,700;0,800;0,900;1,400;1,500;1,600;1,700;1,800;1,900&display=swap");
/*font-family: 'Playfair Display', serif;*/

.start-49 {
  left: 49% !important;
}

.top-45 {
  top: 45% !important;
}

input:focus {
  border: 0;
  outline: none;
}

.search-input {
  background-color: rgba(0, 0, 0, 0);
  border: 0;
  color: #50493e;
  font-family: "Courier Prime", monospace;
  font-size: 2rem;
}

.search-input::-webkit-input-placeholder {
  /* Chrome/Opera/Safari */
  color: #8a8070;
  font-family: "Playfair Display", serif;
}
.search-input::-moz-placeholder {
  /* Firefox 19+ */
  color: #8a8070;
  font-family: "Playfair Display", serif;
}
.search-input:-ms-input-placeholder {
  /* IE 10+ */
  color: #8a8070;
  font-family: "Playfair Display", serif;
}
.search-input:-moz-placeholder {
  /* Firefox 18- */
  color: #8a8070;
  font-family: "Playfair Display", serif;
}
</style>
