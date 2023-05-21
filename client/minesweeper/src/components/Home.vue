<script setup lang='ts'>
import Board from './Board.vue'
import Panel from './Panel.vue'
import GameLayout from './GameLayout.vue'
import TopMenu from './TopMenu.vue'
import Modal from './Modal.vue'

import { ref } from 'vue'
import { addSocketEventHandler } from '@/socket'
import { GAMESETTING } from '@/config'
import boardStore from '@/store/boardStore'
import { activeStore } from '@/store'

const turnCount = ref(1)
const roomId = ref(-1)
const currGameState = ref(GAMESETTING.NEW)
const isPlayerTurn = ref(false)
const modalContent = ref('createOrJoin')

addSocketEventHandler('turnPassed', (data:{ count: number }) => {
    const { count } = data
    turnCount.value = count
    isPlayerTurn.value = !isPlayerTurn.value
})

addSocketEventHandler('roomCreated', (data:{ roomId: number, isPlayerTurn: boolean }) => {
    const { roomId: id, isPlayerTurn: isTurn } = data
    roomId.value = id
    turnCount.value = 1
    currGameState.value = GAMESETTING.WAITING
    isPlayerTurn.value = isTurn
    boardStore.resetBoard()
})

addSocketEventHandler('roomJoined', (data: { isPlayerTurn: boolean }) => {
    const { isPlayerTurn: isTurn } = data
    currGameState.value = GAMESETTING.PLAYING
    isPlayerTurn.value = isTurn
})

addSocketEventHandler('newPlayer', (data: {  }) => {
    currGameState.value = GAMESETTING.PLAYING
})

const displayModal = (v: 'create' | 'join' | 'createOrJoin') => {
    modalContent.value = v
    activeStore.active = false
}
</script>
    
<template>
    <Modal v-if='!activeStore.active' :content='modalContent !== `` ? modalContent : undefined'/>
    <div class='app-container'>
        <GameLayout>
            <template #header>
                <TopMenu :roomId='roomId' :displayModal='displayModal'/>
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
