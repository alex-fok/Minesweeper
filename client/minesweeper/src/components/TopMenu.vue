<script setup lang='ts'>
const props = defineProps({
    roomId: Number,
    gameState: Number,
    displayModal: {
        type: Function,
        default: (v:string) => {}
    }
})

const createRoom = () => {
    props.displayModal('create')
}

const joinRoom = () => {
    props.displayModal('join')
}
</script>
<template>
    <div class='header-container'>
        <div class='menu-item'><div class='room'>Room #{{ roomId !== undefined && roomId >= 0 ? roomId : ' (Loading...)' }}</div></div>
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
