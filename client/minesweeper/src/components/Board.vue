<script setup lang='ts'>
import { socket, addSocketEventHandler } from '@/socket'
import { boardStore } from '@/store'
import { boardSetting } from '@/config'

type blockInfo = {
    x: number,
    y: number,
    bType: number,
    value: number
}
type block = {
    x: number,
    y: number,
    show: string
}

const [BLANK, BOMB, NUMBER] = [0, 1, 2]

const reveal = (i: number) => {
    const y = Math.floor(i / boardSetting.SIZE)
    const x = i % boardSetting.SIZE
    socket.send(JSON.stringify({
        name: 'reveal',
        content: JSON.stringify({x, y})
    }))
}

const modifyBoard = (board:block[], x:number, y:number, show: string) => {
    boardStore.board = board.map((v, _) => {
        if (v.x === x && v.y === y)
            v.show = show
        return v
    })
}

const getDisplayVal = (block: blockInfo) : string => {
    if (block['bType'] === NUMBER) return block['value'].toString()
    return block['bType'] === BOMB ? 'BO' : 'BL'
}

addSocketEventHandler('reveal', (data: {blocks:blockInfo[]}) => {
    const { blocks } = data
    blocks.forEach(block => {
        modifyBoard(boardStore.board, block['x'], block['y'], getDisplayVal(block))
    })
})
</script>
<template>
    <div class='board'>
       <div
          v-for='(block, i) in boardStore.board'
          @click='reveal(i)'
          :key='i'
        >
        {{ block.show }}
        </div>
    </div>
</template>
<style scoped>
    .board {
        justify-content: center;
        column-gap: 1px;
        row-gap:1px;
        flex-grow: 1;
        display: grid;
        grid-template-columns:auto auto auto auto auto auto auto auto auto auto auto 
        auto auto auto auto auto auto auto auto auto auto auto auto auto auto auto;
    }
    .board > div {
        font-family:'Courier New', Courier, monospace;
        background-color:#343434;   
        width: 3vh;
        height: 3vh;
        user-select: none;
    }
    
</style>
