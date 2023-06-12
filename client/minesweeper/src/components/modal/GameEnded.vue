<script setup lang='ts'>
import { gameState } from '@/store';
defineProps({
    close: {
        type: Function,
        default: () => {}
    }
})
const isPlayer = gameState.players[gameState.id] !== undefined
const isWon = isPlayer && gameState.winner === gameState.id
</script>
<template>
    <template v-if='isPlayer'>
        <div class='modal-item'>
            <div class='modal-end-game'>{{ isWon ? 'You Won!' : 'You Lost!' }}</div>
        </div>
        <div class='modal-item'>
            <span class='btn' @click=''>REMATCH?</span>
        </div>
        <div class='modal-close' @click='close()'>&#10005;</div>
    </template>
    <template v-else>
        <div class='modal-item'>
            {{ gameState.players[gameState.winner].alias }} Won!
        </div>
        <div class='modal-close' @click='close()'>&#10005;</div>
    </template>
</template>
<style scoped>
    @import '@/assets/modal.css';
    .modal-end-game {
        font-size: 1.7rem;
    }
</style>

