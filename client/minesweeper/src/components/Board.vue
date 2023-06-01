<script setup lang='ts'>
import { socket, addSocketEventHandler } from '@/socket'
import { gameState } from '@/store'
import { BOARDSETTING, GAMESTATUS } from '@/config'

type blockInfo = {
    x: number,
    y: number,
    bType: number,
    value: number
}

const [BLANK, BOMB, NUMBER] = [0, 1, 2]

const reveal = (i: number) => {
    if (gameState.status !== GAMESTATUS.PLAYING) return
    const y = Math.floor(i / BOARDSETTING.SIZE)
    const x = i % BOARDSETTING.SIZE
    socket.send(JSON.stringify({
        name: 'reveal',
        content: JSON.stringify({x, y})
    }))
}

addSocketEventHandler('reveal', (data: {blocks:blockInfo[]}) => {
    const { blocks } = data
    const getDisplayVal = (block: blockInfo) : string => {
        if (block['bType'] === NUMBER) return block['value'].toString()
        return block['bType'] === BOMB ? 'BO' : 'BL'
    }
    blocks.forEach(block => {
        gameState.board[BOARDSETTING.SIZE * block.y + block.x].show = getDisplayVal(block)
    })
})
</script>
<template>
    <div class='board-container'>
        <div v-if='gameState.status === GAMESTATUS.WAITING_JOIN'>
            Waiting for player to join...
        </div>
        <div v-else class='board-wrapper'>
            <div class='board'>
                <div
                    v-for='(block, i) in gameState.board'
                    @click='reveal(i)'
                    :key='i'
                    >
                    {{ block.show }}
                </div>
            </div>
            <div
                v-if='gameState.status === GAMESTATUS.WAITING_TURN'
                class='overlay'
            >
               <div class='overlay-text'>Waiting for opponent...</div>
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
    .board > div {
        font-family:'Courier New', Courier, monospace;
        background-color:#343434;   
        width: 3vh;
        height: 3vh;
        user-select: none;
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
