<script setup lang='ts'>
import socket from '@/socket';
import { gameState } from '@/store';
defineProps({
    close: {
        type: Function,
        default: () => {}
    }
})
const isPlayer = gameState.players[gameState.id] !== undefined
const isWon = isPlayer && gameState.winner === gameState.id
const oppId = Object.getOwnPropertyNames(gameState.players).find(id => id !== gameState.id)
const opponent = oppId ? gameState.players[oppId] : null

const requestRematch = () => {
    socket.send(JSON.stringify({
        name: 'rematch',
        content: JSON.stringify({
            rematch: true
        })
    }))
}
</script>
<template>
    <template v-if='isPlayer'>
        <div class='modal-row'>
            <div class='modal-item grow'>
                <div class='modal-end-game'>{{ isWon ? 'You Won!' : 'You Lost!' }}</div>
            </div>
            <div class='modal-close' @click='close()'>&#10005;</div>
        </div>
        <div class='modal-row'>
                <div class='scoreboard grow'>
                    <div class='score'>
                        <span class='score-text'>{{ gameState.players[gameState.id].score }}</span>
                        {{ gameState.players[gameState.id].alias }}
                    </div>
                    <div class='score-text'> - </div>
                    <div class='score'>
                        <span class='score-text'>{{ opponent?.score }}</span>
                        {{ opponent?.alias || "Anonymous" }}
                    </div>
                </div>
        </div>
        <div class='modal-row'>
            <div class='modal-item'>
                <span class='btn' @click='requestRematch'>REMATCH?</span>
            </div>
        </div>
    </template>
    <template v-else>
        <div class='moda-row'>
            <div class='modal-item'>
                {{ gameState.players[gameState.winner].alias }} Won!
            </div>
            <div class='modal-close' @click='close()'>&#10005;</div>
        </div>
    </template>
</template>
<style scoped>
    @import '@/assets/modal.css';
    .modal-end-game {
        font-size: 1.7rem;
    }
    .scoreboard {
        display: flex;
        flex-direction: row;
        justify-content: space-between;
        column-gap: .5rem;
    }
    .score {
        display: flex;
        flex-direction: column;
        flex-grow: 1;
        text-align: center;
    }
    .score-text {
        font-size: 4rem;
    }
    .grow {
        flex-grow: 1;
    }
</style>
