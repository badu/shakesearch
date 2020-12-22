import App from './App.vue'
import './assets/css/app.css'
import { createApp } from 'vue';
import { createStore } from 'vuex';

const store = createStore({
    state() {
        return {
            searchTerm: '',
            loading: false,
            results: null,
            lastError: null
        };
    },
    actions: {
        clearResults({ state }) {
            state.searchTerm = '';
            state.results = null;
        },
        async fetchResults({ state }, newSearchTerm) {
            state.searchTerm = newSearchTerm;
            state.loading = true;
            state.results = [];

            const response = await fetch(`http://localhost:8080/search?q=${newSearchTerm}`);

            if (response.status !== 200) {
                state.lastError = response.statusText;
                return
            }

            const reader = response.body.getReader();
            const decoder = new TextDecoder('utf-8');
            const parentState = state;
            reader.read().then(
                ({ value: bytes, done }) => {
                    state.loading = false;
                    if (done) {
                        console.log(`stream closed. no results!`);
                        return
                    }
                    const chunk = decoder.decode(bytes);
                    const entries = chunk.split('\n');
                    entries.slice(0, -1).forEach(payload => {
                        try {
                            const row = JSON.parse(payload);
                            parentState.results.push(row);
                        } catch (err) {
                            console.error(`error parsing json ${err} : ${payload}`);
                        }
                    });

                },
                e => console.error(`The stream has an error and cannot be read from : ${e}`)
            ).catch(err => console.error(`something horrible going on : ${err}`));
        }
    }
});
const app = createApp(App);

app.use(store);

app.mount('#app');