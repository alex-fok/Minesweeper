<script setup lang='ts'>
import { GAMESTATUS } from '@/config'
import { gameState, uiState } from '@/store'
import { computed, ref } from 'vue'

defineProps({
    close: {
        type: Function,
        default: () => {}
    }
})

const page = ref(1)
const maxPage = 2

const prevBtn = computed(() => page.value === 1 && gameState.status !== GAMESTATUS.NEW ? 'btn hidden' : 'btn')
const nextBtn = computed(() => page.value === maxPage ? 'btn hidden' : 'btn')

const nextPage = () => {
    if (page.value < maxPage)
        page.value += 1
}

const previousPage = () => {
    page.value <= 1 ?
        uiState.modal.displayContent('createOrJoin') :
        page.value -= 1
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
        <button :class='prevBtn' @click='previousPage()'>
            {{gameState.status === GAMESTATUS.NEW ? 'BACK' : '&#8592; PREV' }}
        </button>
    </div>
    <div class='modal-item'>
        <button :class='nextBtn' @click='nextPage()'>NEXT &#8594;</button>
    </div>
</div>
</template>
<style></style>
