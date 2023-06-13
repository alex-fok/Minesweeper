<script setup lang='ts'>
import socket from '@/socket'
import Block from './Block.vue'
import { gameState } from '@/store'
import { BOARDSETTING, GAMESTATUS } from '@/config'

const reveal = (i: number) => {
    if (gameState.status !== GAMESTATUS.IN_GAME) return

    const y = Math.floor(i / BOARDSETTING.SIZE)
    const x = i % BOARDSETTING.SIZE
    socket.send(JSON.stringify({
        name: 'reveal',
        content: JSON.stringify({x, y})
    }))
}
</script>
<template>
    <div class='board-container'>
        <div v-if='gameState.status === GAMESTATUS.WAITING_JOIN'>
            Waiting for player to join...
        </div>
        <div v-else class='board-wrapper'>
            <div class='board'>
                <Block
                    v-for='(block, i) in gameState.board'
                    :key='i'
                    :reveal='() => { reveal(i) }'
                    :show='block.show'
                    :owner='block.owner'
                />
            </div>
        </div>
    </div>
</template>
<style scoped>
    .board-container {
        flex-grow: 1;
        height: 100%;
        display: flex;
        align-items:center;
        justify-content: center;
    }
    .board-wrapper {
        position:relative
    }
    .board {
        display: grid;
        column-gap: 1px;
        row-gap: 1px;
        grid-template-columns: repeat(26, auto);
    }
    .overlay {
        position:absolute;
        inset: 0;
        background-color:rgba(52, 52, 52, .7);
    }
    .overlay-text {
        position: absolute;
        top: 50%;
        left: 50%;
        transform: translate(-50%, -50%);
    }
</style>
