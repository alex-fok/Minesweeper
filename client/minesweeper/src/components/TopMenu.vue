<script setup lang='ts'>
import { gameState, uiState } from '@/store'

const createRoom = () => {
    uiState.modal.displayContent('create')
}

const joinRoom = () => {
    uiState.modal.displayContent('join')
}
</script>
<template>
    <div class='header-container'>
        <div class='menu-item'>
            <div class='room'>
                Room #{{ gameState.roomId !== undefined && gameState.roomId >= 0 ? gameState.roomId : ' (Loading...)' }}
            </div>
        </div>
        <div class='menu-item'>
            <div class='button' @click='createRoom()'>
                Create Room
            </div>
        </div>
        <div class='menu-item'>
            <div class='button' @click='joinRoom()'>
                Join Room
            </div>
        </div>
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
        display: flex;
        padding: 0rem 1rem;
        height: 4rem;
        align-items: center;
    }
    .room {
        font-size: 1.2rem;
        user-select: text;
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
