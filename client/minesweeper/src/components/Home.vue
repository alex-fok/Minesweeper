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

addSocketEventHandler('gameStarted', (data: { isPlayerTurn: boolean, id: number }) => {
    gameState.resetBoard()
    const { isPlayerTurn, id } = data
    gameState.status = isPlayerTurn ? PLAYING : WAITING_TURN
    roomId.value = id
})

const displayModal = (v: 'create' | 'join' | 'createOrJoin') => {
    modalContent.value = v
    uiState.active = false
}
</script>
    
<template>
    <Modal v-if='!uiState.active' :content='modalContent !== `` ? modalContent : undefined'/>
    <div class='app-container'>
        <GameLayout>
            <template #header>
                <TopMenu :roomId='roomId' :displayModal='displayModal'/>
            </template>
            <template #default v-if='gameState.status !== NEW'>
                <Board/>
                <Panel :turnCount='turnCount'/> 
            </template>
        </GameLayout>
    </div>
</template>
<style scoped>
    .app-container {
        display: flex;
        flex-direction: column;
        height: 100vh;
    }
</style>
