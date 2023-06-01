<script setup lang='ts'>
import { GAMESTATUS } from '@/config';
import { gameState } from '@/store';
defineProps({
    close: {
        type: Function,
        default: () => {}
    }
})
const isWon = gameState.status === GAMESTATUS.GAME_WON
</script>
<template>
    <div class='modal-item'>
        <div class='modal-end-game'>{{ isWon ? 'You Won!' : 'You Lost!' }}</div>
    </div>
    <div class='modal-close' @click='close()'>&#10005;</div>
</template>
<style scoped>
    @import '@/assets/modal.css';
    .modal-end-game {
        font-size: 1.7rem;
    }
</style>

