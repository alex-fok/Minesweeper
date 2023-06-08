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

const displayModal = (v: 'create' | 'join' | 'createOrJoin' | 'gameEnded') => {
    modalContent.value = v
    uiState.active = false
}

const updateGameStat = (playerScore: number, opponentScore: number, bombsLeft: number) => {
    gameState.player.score = playerScore
    gameState.opponent.score = opponentScore
    gameState.bombsLeft = bombsLeft
}

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

    // Change search query
    const url = new URL(window.location.href)
    url.searchParams.set('room', id.toString())
    history.replaceState({}, '', url)
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
    
    // Change search query
    const url = new URL(window.location.href)
    url.searchParams.set('room', id.toString())
    history.replaceState({}, '', url)
})

addSocketEventHandler('scoreUpdated', (data: { player: number, opponent: number, bombsLeft: number }) => {
    const { player, opponent, bombsLeft } = data
    updateGameStat(player, opponent, bombsLeft)
})

addSocketEventHandler('gameWon', (data: { player: number, opponent: number, bombsLeft: number}) => {
    const { player, opponent, bombsLeft } = data
    updateGameStat(player, opponent, bombsLeft)
    gameState.status = GAMESTATUS.GAME_WON
    displayModal('gameEnded')
})

addSocketEventHandler('gameLost', (data: { player: number, opponent: number, bombsLeft: number}) => {
    const { player, opponent, bombsLeft } = data
    updateGameStat(player, opponent, bombsLeft)
    gameState.status = GAMESTATUS.GAME_LOST
    displayModal('gameEnded')
})

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
