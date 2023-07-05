<script setup lang='ts'>
import { GAMESTATUS } from '@/config'
import { gameState, uiState } from '@/store'
import { ref } from 'vue'
const props = defineProps({
    close: {
        type: Function,
        default: () => {}
    }
})
const page = ref(1)
const maxPage = 2
const nextPage = () => {
    if (page.value < maxPage)
        page.value += 1
}

const previousPage = () => {
    if (page.value <= 1) {
        gameState.status === GAMESTATUS.NEW ?
            uiState.modal.displayContent('createOrJoin') :
            props.close()
    } else {
        page.value -= 1
    }
}
</script>
<template>
<div class='modal-row'>
    <div class='modal-item'>
        <div v-if='page === 1'>
            <p>Placeholder 1</p>
        </div>
        <div v-else>
            <p>Placeholder 2</p>
        </div>
        <div class='modal-close' @click='close()'>&#10005;</div>
    </div>
    
</div>
<div class='modal-row'>
    <div class='modal-item'>
        <button class='btn' @click='previousPage()'>&#8592; PREV</button>
    </div>
    <div class='modal-item'>
        <button class='btn' @click='nextPage()'>NEXT &#8594;</button>
    </div>
</div>
</template>
<style></style>
