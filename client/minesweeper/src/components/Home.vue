<script setup lang='ts'>
    import Board from './Board.vue'
    import Panel from './Panel.vue'
    import GameLayout from './GameLayout.vue'
    import TopMenu from './TopMenu.vue'
    
    import { ref } from 'vue'
    import { addSocketEventHandler } from '@/socket'
    import { gameState } from '@/config'
    import { store } from '@/store/boardStore'
    
    const container = ref('app-container')
    const turnCount = ref(1)
    const roomId = ref(-1)
    const currGameState = ref(gameState.NEW)
    const isPlayerTurn = ref(false)
    
    addSocketEventHandler('turnPassed', (data:{ count: number }) => {
        const { count } = data
        turnCount.value = count
        isPlayerTurn.value = !isPlayerTurn.value
    })

    addSocketEventHandler('gameCreated', (data:{ roomId: number, isPlayerTurn: boolean }) => {
        const { roomId: id, isPlayerTurn: isTurn } = data
        roomId.value = id
        turnCount.value = 1
        currGameState.value = gameState.WAITING
        isPlayerTurn.value = isTurn
        store.resetBoard()
    })

    addSocketEventHandler('gameJoined', (data: { isPlayerTurn: boolean }) => {
        const { isPlayerTurn: isTurn } = data
        currGameState.value = gameState.PLAYING
        isPlayerTurn.value = isTurn
    })

</script>
    
<template>
    <div :class='container'>
        <TopMenu :roomId='roomId'/>
        <GameLayout>
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
