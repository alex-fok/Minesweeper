<script setup lang='ts'>
import { gameState, uiState } from '@/store'
import { ref } from 'vue'
defineProps({
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
        uiState.modal.displayContent('createOrJoin')
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
        <span class='btn' @click='previousPage()'>&#8592; PREV</span>
    </div>
    <div class='modal-item'>
        <span class='btn' @click='nextPage()'>NEXT &#8594;</span>
    </div>
</div>
</template>
<style></style>
