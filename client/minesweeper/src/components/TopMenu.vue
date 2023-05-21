<script setup lang='ts'>
import { socket } from '@/socket'
import { activeStore } from '@/store'


defineProps({
    roomId: Number,
    gameState: Number
})

const createRoom = () => {
    socket.send(JSON.stringify({name: 'newRoom'}))
}

const joinRoom = () => {
    activeStore.updateActivity(false)
}
</script>
<template>
    <div class='header-container'>
        <div class='menu-item'><div class='room'>Room #{{ roomId && roomId >= 0 ? roomId : ' (Loading...)' }}</div></div>
        <div class='menu-item'><div class='button' @click='createRoom()'>Create Room</div></div>
        <div class='menu-item'><div class='button' @click='joinRoom()'>Join Room</div></div>
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
