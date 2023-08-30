<script setup lang='ts'>
import { GAMESTATUS } from '@/config';
import { gameState, uiState } from '@/store';
import { computed, ref } from 'vue';
import Edit from './icon/edit.vue';

const [hovering, setHovering] = [ref(''), (id: string) => { hovering.value = id }]

const isGameStarted = computed(() => 
    gameState.status !== GAMESTATUS.NEW &&
    gameState.status !== GAMESTATUS.WAITING_JOIN
)
const getPlayerColor = (id: string) => {
    const color = uiState.playerColor[id] || '0, 0, 0'
    const isHovering = id === hovering.value
    const opacity = gameState.players[id].isTurn ? 1 : isHovering ? .4 : .2
    return `rgba(${color}, ${opacity})`
}
const playerStyle = computed(() => 
    (id: string) => {
        const playerColor = getPlayerColor(id)
        return ({
            border: `.1rem solid ${playerColor}`,
            color:`${playerColor}`
        })
    }
)

const editName = (id: string) => {
    if (id !== gameState.id) return
    uiState.modal.displayContent('playerAlias')
}
</script>
<template>
    <div v-if='isGameStarted' class='side-container'>
        <div v-if='gameState.capacity > 2' class='player-container'>
            <div
                v-for='player in gameState.players'
                class='player-expand-item'
                :style='playerStyle(player.id)'
            >
                <div class='player'
                    :onmouseenter='() => { setHovering(player.id) }'
                    :onmouseleave='() => { setHovering(``) }'
                >
                    <span :class='player.isTurn ? `` : `hidden`'>></span>
                    <span
                        :class='`${gameState.id === player.id ? `name self` : `name`}`'
                        :onclick='() => { editName(player.id) }'
                    >
                        {{ player.alias }}
                        {{  gameState.id === player.id ? '(You)' : '' }}
                        {{ !player.isOnline ? '(Offline)': '' }}
                        <Edit v-if='gameState.id === player.id && hovering === player.id'
                            :fill='getPlayerColor(player.id)'
                            size='1rem'
                        />
                    </span>
                </div>
                <div class='score'>{{ player.score }}</div>
            </div>
        </div>
        <div v-else class='player-container'>
            <div
                v-for='player in gameState.players'
                class='player-collapse-item'
                :onmouseenter='() => { setHovering(player.id) }'
                :onmouseleave='() => { setHovering(``) }'
                :style='playerStyle(player.id)'
            >
                <div class='player'>
                    <span :class='player.isTurn ? `` : `hidden`'>></span>
                    <span
                        :class='`${gameState.id === player.id ? `name self` : `name`}`'
                        :onclick='()=> { editName(player.id) }'
                    >
                        {{ player.alias }}
                        {{ gameState.id === player.id ? '(You)' : '' }}
                        {{ !player.isOnline ? '(Offline)' : '' }}
                        <Edit v-if='gameState.id === player.id && hovering === player.id'
                            :fill='getPlayerColor(player.id)'
                            size='1rem'
                        />
                    </span>
                </div>
                <div class='score'>{{ player.score }}</div>
            </div>
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
    .player {
        display: grid;
        grid-template-columns: 2rem 1fr;
        word-break: break-all;
    }
    .name {
        background: transparent;
        color: inherit;
        border: 0;
        outline: 0;
        width: auto;
    }
    .name.self:hover {
        text-decoration: underline;
        cursor:pointer;
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
