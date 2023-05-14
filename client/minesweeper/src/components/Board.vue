<script setup lang='ts'>
    import { socket, addSocketEventHandler } from '@/socket'
    import { store }from '@/store/boardStore'
    const reveal = (i: number) => {
        const y = Math.floor(i / 26)
        const x = i % 26
        socket.send(JSON.stringify({
            name: "reveal",
            content: JSON.stringify({x, y})
        }))
    }
</script>
<template>
    <div class='container'>
       <div
          v-for='(block, i) in store.board'
          @click='reveal(i)'
          :key='i'
        >
        {{ block.show }}
        </div>
    </div>
</template>
<script lang='ts'>
    const [BLANK, BOMB, NUMBER] = [0, 1, 2]
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

    store.board = store.board.map((_, i) => ({
        x: i % 26,
        y: Math.floor(i / 26),
        show: ""
    }))
    const modifyBoard = (board:block[], x:number, y:number, show: string) => {
        store.board = board.map((v, _) => {
            if (v.x === x && v.y === y)
                v.show = show
            return v
        })
    }
    const getDisplayVal = (block: blockInfo) : string => {
        if (block["bType"] === NUMBER) return block["value"].toString()
        return block["bType"] === BOMB ? "BO" : "BL"
    }
    addSocketEventHandler("reveal", (event: any) => {
        const blocks:blockInfo[] = JSON.parse(JSON.parse(event.data))
        blocks.forEach(block => {
            modifyBoard(store.board, block["x"], block["y"], getDisplayVal(block))
        })
    })
    export default {
        name: 'Board',
    }
</script>
<style scoped>
    .container {
        width:fit-content;
        height:fit-content;
        column-gap: 1px;
        row-gap:1px;
        display: grid;
        grid-template-columns:auto auto auto auto auto auto auto auto auto auto auto 
        auto auto auto auto auto auto auto auto auto auto auto auto auto auto auto;
    }
    .container > div {
        background-color:lightblue;   
        width: 3vh;
        height: 3vh;
    }
    
</style>
