<script setup lang='ts'>
    import { ref } from 'vue'
    import { socket } from '@/socket'
    import config from '@/config'

    const container = ref('header-container')
    const room = ref('room')
    const menuItem = ref('menu-item')
    const button = ref('button')

    const props = defineProps({
        roomId: Number,
        gameState: Number
    })

    const createGame = () => {
        socket.send(JSON.stringify({name: 'newGame'}))
    }

    const joinGame = () => {
        if (props.gameState !== config.gameState.NEW) return
    }
</script>
<template>
    <div :class='container'>
        <div :class='menuItem'><div :class='room'>Room #{{ roomId && roomId >= 0 ? roomId : ' (Loading...)' }}</div></div>
        <div :class='menuItem'><div :class='button' @click='createGame'>Create Game</div></div>
        <div :class='menuItem'><div :class='button' @click='joinGame'>Join Game</div></div>
    </div>
</template>
<style scoped>
    .header-container {
        border-bottom: .1rem solid #9F9F9F;
        padding: 0rem 1rem;
        display: flex;
        flex-direction: row;
        user-select: none;
    }
    .menu-item {
        position: relative;
        padding: 0rem 1rem;
        line-height: 4rem;
        vertical-align: middle;
    }
    .room {
        font-size: 1.2rem;
    }
    .button {
        color:#9F9F9F;
        font-size: .8rem;
        cursor: pointer;
    }
    .button:hover {
        color: white;
    }
</style>
