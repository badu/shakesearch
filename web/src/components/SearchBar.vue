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
        <div class="col input">
          <div class="text">
            <input
              type="text"
              class="search-input"
              :value="searchTerm"
              v-on:keyup.enter="performSearch"
              placeholder="What art thee looking f'r?"
            />
            <button class="clear">
              <svg viewBox="0 0 24 24">
                <path class="line" d="M2 2L22 22" />
                <path class="long" d="M9 15L20 4" />
                <path class="arrow" d="M13 11V7" />
                <path class="arrow" d="M17 11H13" />
              </svg>
            </button>
          </div>
        </div>
      </div>
      <img src="../assets/svg/signature.svg" alt="sketch" />
    </div>
  </div>
</template>

<script>
// eslint-disable-next-line no-undef
const { to, set } = gsap;

function delay(fn, ms) {
  let timer = 0;
  return function (...args) {
    clearTimeout(timer);
    timer = setTimeout(fn.bind(this, ...args), ms || 0);
  };
}

function getPoint(point, i, a, smoothing) {
  let cp = (current, previous, next, reverse) => {
      let p = previous || current,
        n = next || current,
        o = {
          length: Math.sqrt(
            Math.pow(n[0] - p[0], 2) + Math.pow(n[1] - p[1], 2)
          ),
          angle: Math.atan2(n[1] - p[1], n[0] - p[0]),
        },
        angle = o.angle + (reverse ? Math.PI : 0),
        length = o.length * smoothing;
      return [
        current[0] + Math.cos(angle) * length,
        current[1] + Math.sin(angle) * length,
      ];
    },
    cps = cp(a[i - 1], a[i - 2], point, false),
    cpe = cp(point, a[i - 1], a[i + 1], true);
  return `C ${cps[0]},${cps[1]} ${cpe[0]},${cpe[1]} ${point[0]},${point[1]}`;
}

function getPath(x, smoothing) {
  return [
    [2, 2],
    [12 - x, 12 + x],
    [22, 22],
  ].reduce(
    (acc, point, i, a) =>
      i === 0
        ? `M ${point[0]},${point[1]}`
        : `${acc} ${getPoint(point, i, a, smoothing)}`,
    ""
  );
}

export default {
  props: ["searchTerm"],
  emits: ["search", "clear-search"],
  methods: {
    performSearch(event) {
      if (event.target.value === "") {
        event.target.placeholder = "Art thee looking f'r something?";
        return;
      }
      this.$emit("search", event.target.value);
    },
  },
  mounted() {
    this.$nextTick(function () {
      // Code that will run only after the entire view has been rendered
      document.querySelectorAll(".input").forEach((elem) => {
        let clear = elem.querySelector(".clear"),
          input = elem.querySelector("input"),
          { classList } = elem,
          svgLine = clear.querySelector(".line"),
          svgLineProxy = new Proxy(
            {
              x: null,
            },
            {
              set(target, key, value) {
                target[key] = value;
                if (target.x !== null) {
                  svgLine.setAttribute("d", getPath(target.x, 0.1925));
                }
                return true;
              },
              get(target, key) {
                return target[key];
              },
            }
          );

        svgLineProxy.x = 0;

        input.addEventListener(
          "input",
          delay(() => {
            let bool = input.value.length;

            to(elem, {
              "--clear-scale": bool ? 1 : 0,
              duration: bool ? 0.5 : 0.15,
              ease: bool ? "elastic.out(1, .7)" : "none",
            });
            to(elem, {
              "--clear-opacity": bool ? 1 : 0,
              duration: 0.15,
            });
          }, 250)
        );

        clear.addEventListener("click", () => {
          classList.add("clearing");
          set(elem, {
            "--clear-swipe-left": (input.offsetWidth - 16) * -1 + "px",
          });
          to(elem, {
            keyframes: [
              {
                "--clear-rotate": "45deg",
                duration: 0.25,
              },
              {
                "--clear-arrow-x": "2px",
                "--clear-arrow-y": "-2px",
                duration: 0.15,
              },
              {
                "--clear-arrow-x": "-3px",
                "--clear-arrow-y": "3px",
                "--clear-swipe": "-3px",
                duration: 0.15,
                onStart() {
                  to(svgLineProxy, {
                    x: 3,
                    duration: 0.1,
                    delay: 0.05,
                  });
                },
              },
              {
                "--clear-swipe-x": 1,
                "--clear-x": input.offsetWidth * -1 + "px",
                duration: 0.45,
                onComplete() {
                  input.value = "";
                  input.focus();
                  to(elem, {
                    "--clear-arrow-offset": "4px",
                    "--clear-arrow-offset-second": "4px",
                    "--clear-line-array": "8.5px",
                    "--clear-line-offset": "27px",
                    "--clear-long-offset": "24px",
                    "--clear-rotate": "0deg",
                    "--clear-arrow-o": 1,
                    duration: 0,
                    delay: 0.7,
                    onStart() {
                      classList.remove("clearing");
                    },
                  });
                  to(elem, {
                    "--clear-opacity": 0,
                    duration: 0.2,
                    delay: 0.55,
                  });
                  to(elem, {
                    "--clear-arrow-o": 0,
                    "--clear-arrow-x": "0px",
                    "--clear-arrow-y": "0px",
                    "--clear-swipe": "0px",
                    duration: 0.15,
                  });
                  to(svgLineProxy, {
                    x: 0,
                    duration: 0.45,
                    ease: "elastic.out(1, .75)",
                  });
                },
              },
              {
                "--clear-swipe-x": 0,
                "--clear-x": "0px",
                duration: 0.4,
                delay: 0.35,
              },
            ],
          });
          to(elem, {
            "--clear-arrow-offset": "0px",
            "--clear-arrow-offset-second": "8px",
            "--clear-line-array": "28.5px",
            "--clear-line-offset": "57px",
            "--clear-long-offset": "17px",
            duration: 0.2,
          });

          const dispatcher = this.$emit;
          setTimeout(() => dispatcher("clear-search"), 1000);

        });
      });
    });
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

.input {
  --placeholder-color: #c9c9d9;
  --placeholder-color-hover: #babac9;
  --close: #818190;
  --close-light: #babac9;
  --close-background: #f1f1fa;
  display: -webkit-box;
  display: flex;
  -webkit-box-align: center;
  align-items: center;
  position: relative;
  background: var(--background);
  -webkit-transition: box-shadow 0.2s;
  transition: box-shadow 0.2s;
  --clear-x: 0px;
  --clear-swipe-left: 0px;
  --clear-swipe-x: 0;
  --clear-swipe: 0px;
  --clear-scale: 0;
  --clear-rotate: 0deg;
  --clear-opacity: 0;
  --clear-arrow-o: 1;
  --clear-arrow-x: 0px;
  --clear-arrow-y: 0px;
  --clear-arrow-offset: 4px;
  --clear-arrow-offset-second: 4px;
  --clear-line-array: 8.5px;
  --clear-line-offset: 27px;
  --clear-long-array: 8.5px;
  --clear-long-offset: 24px;
}

.input.clearing,
.input:focus-within {
  --border-width: 1.5px;
  --border: var(--border-active);
  --shadow: var(--shadow-active);
}
.input.clearing {
  --close-background: transparent;
  --clear-arrow-stroke: var(--close-light);
}

.input .clear {
  -webkit-appearance: none;
  position: relative;
  outline: none;
  z-index: 1;
  padding: 0;
  margin: 12px 12px 12px 0;
  border: none;
  background: var(--b, transparent);
  -webkit-transition: background 0.2s;
  transition: background 0.2s;
  border-radius: 50%;
  opacity: var(--clear-opacity);
  -webkit-transform: scale(var(--clear-scale)) translateZ(0);
  transform: scale(var(--clear-scale)) translateZ(0);
}
.input .clear:before {
  content: "";
  position: absolute;
  top: 0;
  bottom: 0;
  right: 12px;
  left: var(--clear-swipe-left);
  background: var(--background);
  -webkit-transform-origin: 100% 50%;
  transform-origin: 100% 50%;
  -webkit-transform: translateX(var(--clear-swipe)) scaleX(var(--clear-swipe-x))
    translateZ(0);
  transform: translateX(var(--clear-swipe)) scaleX(var(--clear-swipe-x))
    translateZ(0);
}
.input .clear svg {
  display: block;
  position: relative;
  z-index: 1;
  width: 24px;
  height: 24px;
  outline: none;
  cursor: pointer;
  fill: none;
  stroke-width: 1.5;
  stroke-linecap: round;
  stroke-linejoin: round;
  stroke: var(--close);
  -webkit-transform: translateX(var(--clear-x)) rotate(var(--clear-rotate))
    translateZ(0);
  transform: translateX(var(--clear-x)) rotate(var(--clear-rotate))
    translateZ(0);
}
.input .clear svg path {
  -webkit-transition: stroke 0.2s;
  transition: stroke 0.2s;
}
.input .clear svg path.arrow {
  stroke: var(--clear-arrow-stroke, var(--close));
  stroke-dasharray: 4px;
  stroke-dashoffset: var(--clear-arrow-offset);
  opacity: var(--clear-arrow-o);
  -webkit-transform: translate(var(--clear-arrow-x), var(--clear-arrow-y))
    translateZ(0);
  transform: translate(var(--clear-arrow-x), var(--clear-arrow-y)) translateZ(0);
}
.input .clear svg path.arrow:last-child {
  stroke-dashoffset: var(--clear-arrow-offset-second);
}
.input .clear svg path.line {
  stroke-dasharray: var(--clear-line-array) 28.5px;
  stroke-dashoffset: var(--clear-line-offset);
}
.input .clear svg path.long {
  stroke: var(--clear-arrow-stroke, var(--close));
  stroke-dasharray: var(--clear-long-array) 15.5px;
  stroke-dashoffset: var(--clear-long-offset);
  opacity: var(--clear-arrow-o);
  -webkit-transform: translate(var(--clear-arrow-x), var(--clear-arrow-y))
    translateZ(0);
  transform: translate(var(--clear-arrow-x), var(--clear-arrow-y)) translateZ(0);
}
.input .clear:hover {
  --b: var(--close-background);
}
</style>
