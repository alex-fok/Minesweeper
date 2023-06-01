<script setup lang='ts'>
import Board from './Board.vue'
import Panel from './Panel.vue'
import GameLayout from './GameLayout.vue'
import TopMenu from './TopMenu.vue'
import Modal from './Modal.vue'

import { ref } from 'vue'
import { addSocketEventHandler } from '@/socket'
import { GAMESTATUS } from '@/config'
import { gameState, uiState } from '@/store'

const { NEW, WAITING_JOIN, PLAYING, WAITING_TURN } = GAMESTATUS

const turnCount = ref(1)
const roomId = ref(-1)
const modalContent = ref('createOrJoin')

addSocketEventHandler('turnPassed', (data:{ count: number }) => {
    const { count } = data
    turnCount.value = count
    gameState.status = gameState.status === PLAYING ? WAITING_TURN : PLAYING
})

addSocketEventHandler('roomCreated', (data:{ roomId: number }) => {
    const { roomId: id } = data
    roomId.value = id
    turnCount.value = 1
    gameState.status = WAITING_JOIN
    gameState.resetBoard()
})

addSocketEventHandler('gameStarted', (data: {
    isPlayerTurn: boolean,
    id: number,
    bombsLeft: number,
    player: { alias: string, score: number },
    opponent: { alias: string, score: number }
}) => {
    gameState.resetBoard()
    const { isPlayerTurn, id, bombsLeft, player, opponent } = data
    gameState.status = isPlayerTurn ? PLAYING : WAITING_TURN
    roomId.value = id
    
    gameState.bombsLeft = bombsLeft
    
    gameState.player = {
        name: player.alias,
        score: player.score 
    }
    gameState.opponent = {
        name: opponent.alias,
        score: opponent.score
    }
})

addSocketEventHandler('scoreUpdated', (data: { player: number, opponent: number, bombsLeft: number }) => {
    const { player, opponent, bombsLeft } = data
    gameState.player.score = player
    gameState.opponent.score = opponent
    gameState.bombsLeft = bombsLeft
})

const displayModal = (v: 'create' | 'join' | 'createOrJoin') => {
    modalContent.value = v
    uiState.active = false
}
</script>
<template>
    <Modal v-if='!uiState.active' :content='modalContent !== `` ? modalContent : undefined'/>
    <GameLayout>
        <template #header>
            <TopMenu :roomId='roomId' :displayModal='displayModal'/>
        </template>
        <template #default v-if='gameState.status !== NEW'>
            <Board/>
            <Panel :turnCount='turnCount'/> 
        </template>
    </GameLayout>
</template>
