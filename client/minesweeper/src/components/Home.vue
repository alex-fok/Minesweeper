<script setup lang='ts'>
    import Board from './Board.vue'
    import Panel from './Panel.vue'
    import GameLayout from './GameLayout.vue'
    import TopMenu from './TopMenu.vue'
    import { ref } from 'vue'
    import { addSocketEventHandler } from '@/socket'
    
    const container = ref('app-container')
    const turnCount = ref(1)
    const roomId = ref(0)
    
    addSocketEventHandler('turn', (data:string) => {
        const { count } = JSON.parse(data)
        turnCount.value = count
    })
    addSocketEventHandler('room_id', (data:string) => {
        const id = JSON.parse(data) 
        roomId.value = id
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
