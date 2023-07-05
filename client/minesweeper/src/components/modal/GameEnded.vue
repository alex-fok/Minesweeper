<script setup lang='ts'>
import socket from '@/socket';
import { gameState } from '@/store';
import { ref } from 'vue';
const props = defineProps({
    close: {
        type: Function,
        default: () => {}
    }
})

const isRematchClicked = ref(false)

const isWon = gameState.isPlayer && gameState.winner === gameState.id
const oppId = Object.keys(gameState.players).find(id => id !== gameState.id)
const opponent = oppId ? gameState.players[oppId] : null

const requestRematch = (rematch: boolean) => {
    isRematchClicked.value = true
    socket.send(JSON.stringify({
        name: 'rematch',
        content: JSON.stringify({
            rematch
        })
    }))
    if (!rematch) props.close()
}
</script>
<template>
    <template v-if='gameState.isPlayer'>
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
                        {{ gameState.players[gameState.id].alias }} (You)
                    </div>
                    <div class='score-text'> - </div>
                    <div class='score'>
                        <span class='score-text'>{{ opponent?.score }}</span>
                        {{ opponent?.alias || "Anonymous" }}
                    </div>
                </div>
        </div>
        <div class='modal-row'>
            <div v-if='isRematchClicked' class='modal-item'>
                Waiting for opponent response...
            </div>
            <template v-else>
                <div class='modal-item'>
                    REMATCH?
                </div>
                <div class='modal-item'>
                    <button class='btn' @click='() => { requestRematch(true) }'>YES</button>
                </div>
                <div class='modal-item'>
                    <button class='btn' @click='() => { requestRematch(false) }'>NO</button>
                </div>
            </template>
        </div>
    </template>
    <template v-else>
        <div class='moda-row'>
            <div class='modal-item'>
                {{ gameState.players[gameState.winner].alias }} Won!
            </div>
            <div class='modal-close' @click='() => { requestRematch(false) }'>&#10005;</div>
        </div>
    </template>
</template>
<style scoped>
    @import '@/assets/styles/modal.css';
    .modal-end-game {
        font-size: 3rem;
        width: 20rem;
    }
    .scoreboard {
        display: flex;
        flex-direction: row;
        justify-content: space-between;
        column-gap: .5rem;
        padding: 0rem 0rem 1rem 0rem;
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
