import { reactive } from 'vue'

export const store = reactive({
    board: Array(26 * 26).fill({})
})
