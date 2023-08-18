<script setup lang='ts'>
import { GAMESTATUS } from '@/config';
import { gameState, uiState } from '@/store';
import { computed } from 'vue';

const isGameStarted = computed(() => 
    gameState.status !== GAMESTATUS.NEW &&
    gameState.status !== GAMESTATUS.WAITING_JOIN
)
const playerStyle = computed(() => {
    return (id: string) => {
        const color = uiState.playerColor[id] || '0, 0, 0'
        const opacity = gameState.players[id].isTurn ? 1 : .2
        return ({
            border: `.1rem solid rgba(${color}, ${opacity})`,
            color:`rgba(${color}, ${opacity})`
        })
    }
})
</script>
<template>
    <div v-if='isGameStarted' class='side-container'>
        <div v-if='gameState.capacity <= 2' class='player-container'>
            <template v-for='player in gameState.players'>
                <div
                    :class='player.isTurn ? `player-expand-item selected` : `player-expand-item`'
                    :style='playerStyle(player.id)'
                >
                    <div class='player-name'>
                        <span :class='player.isTurn ? `` : `hidden`'>>></span>
                        <span>
                            {{ player.alias }}
                            {{ gameState.id === player.id ? '(You)': '' }}
                            {{ !player.isOnline ? '(Offline)': '' }}
                        </span>
                    </div>
                    <div class='score'>{{ player.score }}</div>
                </div>
            </template>
        </div>
        <div v-else class='player-container'>
            <template v-for='player in gameState.players'>
                <div
                    :class='player.isTurn ? `player-collapse-item selected` : `player-collapse-item`'
                    :style='playerStyle(player.id)'
                >
                    <div class='player-name'>
                        <span :class='player.isTurn ? `` : `hidden`'>></span>
                        <span>
                            {{ player.alias }}
                            {{ gameState.id === player.id ? '(You)' : '' }}
                            {{ !player.isOnline ? '(Offline)' : '' }}
                        </span>
                    </div>
                    <div class='score'>{{ player.score }}</div>
                </div>
            </template>
        </div>
        <div>Bombs Left: {{ gameState.bombsLeft }} / {{ gameState.boardConfig.bomb }}</div>
    </div>
</template>
<style scoped>
    .side-container {
        padding: 2rem 2rem;
        box-sizing: border-box;
        background-color: #222;
        width: 25rem;
        height: 100%;
        display: flex;
        flex-direction: column;
        color: white;
        row-gap: 1rem;
    }
    .player-container {
        display: flex;
        flex-direction: column;
        row-gap: 1rem;
        user-select: none;
    }
    .player-expand-item {
        max-width: 13rem;
        height: 13rem;
        border: 1px solid #9F9F9F;
        border-radius: .5rem;
        flex-grow: 1;
        padding: 1rem 1rem;
        display: flex;
        flex-direction: column;
    }
    .player-collapse-item {
        border: 1px solid #9F9F9F;
        border-radius: .5rem;
        padding: .5rem .5rem;
        display: flex;
        justify-content: space-between;
        flex-direction: row;
    }
    .selected {
        border-color: white;
        color: white;
    }
    .player-name {
        display: grid;
        grid-template-columns: 2rem 1fr;
        word-break: break-all;
    }
    .player-expand-item .score {
        font-size: 6rem;
        text-align: center;
        margin: auto;
    }
    .hidden {
        visibility: hidden;
    }
</style>
