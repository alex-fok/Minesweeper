<script setup lang='ts'>
import Board from './Board.vue'
import Panel from './Panel.vue'
import GameLayout from './GameLayout.vue'
import TopMenu from './TopMenu.vue'
import Modal from './Modal.vue'

import { ref } from 'vue'
import { addSocketEventHandler } from '@/socket'
import { gameState } from '@/config'
import boardStore from '@/store/boardStore'
import { activeStore } from '@/store'

const turnCount = ref(1)
const roomId = ref(-1)
const currGameState = ref(gameState.NEW)
const isPlayerTurn = ref(false)

addSocketEventHandler('turnPassed', (data:{ count: number }) => {
    const { count } = data
    turnCount.value = count
    isPlayerTurn.value = !isPlayerTurn.value
})

addSocketEventHandler('roomCreated', (data:{ roomId: number, isPlayerTurn: boolean }) => {
    const { roomId: id, isPlayerTurn: isTurn } = data
    roomId.value = id
    turnCount.value = 1
    currGameState.value = gameState.WAITING
    isPlayerTurn.value = isTurn
    boardStore.resetBoard()
})

addSocketEventHandler('roomJoined', (data: { isPlayerTurn: boolean }) => {
    const { isPlayerTurn: isTurn } = data
    currGameState.value = gameState.PLAYING
    isPlayerTurn.value = isTurn
})

addSocketEventHandler('newPlayer', (data: {  }) => {
    currGameState.value = gameState.PLAYING
})
</script>
    
<template>
    <div class='app-container'>
        <Modal v-if='!activeStore.active' />
        <GameLayout>
            <template #header>
                <TopMenu :roomId='roomId'/>
            </template>
            <Board />
            <Panel :turnCount='turnCount'/> 
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
