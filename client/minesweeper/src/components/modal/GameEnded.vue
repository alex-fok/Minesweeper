<script setup lang='ts'>
import socket from '@/socket';
import { gameState, uiState } from '@/store';
import { computed, ref } from 'vue';
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

const fillStyle = computed(() => {
    return (id: string) => {
        const color = uiState.playerColor[id] || '0, 0, 0'
        const score = gameState.players[id].score || 0
        const opacity = gameState.winner === id ? 1 : .2
        return ({
            width: `${score/gameState.boardConfig.bomb * 100}%`,
            backgroundColor: `rgba(${color}, ${opacity})`
        })
    }
})

const playerStyle = computed(() => {
    return (id: string) => {
        const color = uiState.playerColor[id] || '0, 0, 0'
        const opacity = gameState.winner === id ? 1 : .2
        return ({
            border: `1px solid rgba(${color}, ${opacity})`,
        })
    }
})
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
                    <div v-for='player in gameState.players'
                        class='score-wrapper'
                        :style='playerStyle(player.id)'
                    >
                        <div class='score-fill' :style='fillStyle(player.id)'></div>
                        <div class='score'>
                            <span>{{ player.alias }}{{ gameState.id === player.id ? ' (You)': '' }}</span>
                            <span>{{ player.score }}</span>
                        </div>
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
        flex-direction: column;
        row-gap: .5rem;
        user-select: none;
        color: white;
    }
    .score-wrapper {
        border-radius: .5rem;
        position: relative;
        overflow: hidden;
    }
    .score-fill {
        position: absolute;
        top: 0;
        left: 0;
        height: 100%;
        background-color: rgba(159, 159, 159, .2);
        z-index: -99;
    }
    .score {
        display: flex;
        flex-direction: row;
        justify-content: space-between;
        padding: .5rem;
    }
    .score-text {
        font-size: 4rem;
    }
    .grow {
        flex-grow: 1;
    }
</style>
