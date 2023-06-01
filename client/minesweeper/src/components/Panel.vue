<script setup lang='ts'>
import gameStatus from '@/config/gameStatus';
import { gameState } from '@/store';
import { computed } from 'vue';

const isTurn = computed(() => {
    return gameState.status === gameStatus.PLAYING
})
const isOppTurn = computed(() => {
    return gameState.status === gameStatus.WAITING_TURN
})
const isGameStarted = computed(() => 
    gameState.status !== gameStatus.NEW &&
    gameState.status !== gameStatus.WAITING_JOIN
    
)
</script>
<template>
    <div v-if='isGameStarted' class='side-container'>
        <div class='player-list'>
            <div :class='isTurn ? `player-item selected` : `player-item`'>
                <div class='player-name'>
                    <span :class='isTurn ? `` : `hidden`'>>></span>
                    <span>{{ gameState.player.name }}</span>
                </div>
                <div class='score'>{{ gameState.player.score }}</div>
            </div>
            <div :class='isOppTurn ? `player-item selected` : `player-item`'>
                <div class='player-name'>
                    <span :class='isOppTurn ? `` : `hidden`'>>></span>
                    <span>{{ gameState.opponent.name }}</span>
                </div>
                <div class='score'>{{ gameState.opponent.score }}</div>
            </div>
        </div>
        <div>Left: {{ gameState.bombsLeft }}</div>
    </div>
</template>
<style scoped>
    .side-container {
        padding: 2rem 2rem;
        box-sizing: border-box;
        background-color: #222;
        width: 18rem;
        height: 100%;
        display: flex;
        flex-direction: column;
        color: #9F9F9F;
        row-gap: 1rem;
    }
    .player-list {
        display: flex;
        flex-direction: column;
        row-gap: 1rem;
        height: 70%;
        user-select: none;
    }
    .player-item {
        border: 1px solid #9F9F9F;
        border-radius: .5rem;
        flex-grow: 1;
        padding: 1rem 1rem;
        display: flex;
        flex-direction: column;
    }
    .player-item.selected {
        border-color: white;
        color: white;
    }
    .player-name {
        display: grid;
        grid-template-columns: 2rem 1fr;
        font-size: 1.2rem;
        word-break: break-all;
    }
    .score {
        font-size: 6rem;
        text-align: center;
        margin: auto;
    }
    .hidden {
        visibility: hidden;
    }
</style>
