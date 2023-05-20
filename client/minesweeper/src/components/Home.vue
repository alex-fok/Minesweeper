<script setup lang='ts'>
    import Board from './Board.vue'
    import Panel from './Panel.vue'
    import GameLayout from './GameLayout.vue'
    import TopMenu from './TopMenu.vue'
    
    import { ref } from 'vue'
    import { addSocketEventHandler } from '@/socket'
    import config from '@/config'
    
    const container = ref('app-container')
    const turnCount = ref(1)
    const roomId = ref(0)
    const currState = ref(config.gameState.NEW)
    const isPlayerTurn = ref(false)
    
    addSocketEventHandler('turn', (data:{ count: number }) => {
        const { count } = data
        turnCount.value = count
        isPlayerTurn.value = !isPlayerTurn.value
    })

    addSocketEventHandler('gameCreated', (data:{ roomId: number, isPlayerTurn: boolean }) => {
        const { roomId: id, isPlayerTurn: isPlayable } = data
        roomId.value = id
        turnCount.value = 1
        currState.value = config.gameState.WAITING
        isPlayerTurn.value = isPlayable
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
