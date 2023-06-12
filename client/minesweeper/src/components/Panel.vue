<script setup lang='ts'>
import { GAMESTATUS } from '@/config';
import { gameState } from '@/store';
import { computed } from 'vue';

const isGameStarted = computed(() => 
    gameState.status !== GAMESTATUS.NEW &&
    gameState.status !== GAMESTATUS.WAITING_JOIN
)
</script>
<template>
    <div v-if='isGameStarted' class='side-container'>
        <div class='player-list'>
            <template v-for='player in gameState.players'>
                <div :class='player.isTurn ? `player-item selected` : `player-item`'>
                    <div class='player-name'>
                        <span :class='player.isTurn ? `` : `hidden`'>>></span>
                        <span>{{ player.alias }}</span>
                    </div>
                    <div class='score'>{{ player.score }}</div>
                </div>
            </template>
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
