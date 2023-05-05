<script setup lang='ts'>
    import { socket, addSocketEvent } from '@/socket'
    import { store }from '@/store/boardStore'
    
    const reveal = (i: number) => {
        const y = Math.floor(i / 26)
        const x = i % 26
        socket.send(JSON.stringify({x, y}))
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
    const blockType = ["BLANK", "BOMB", "NUMBER"]
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
        return board.map((v, _) => {
            if (v.x === x && v.y === y)
                v.show = show
            return v
        })
    }
    addSocketEvent("reveal", (event: any) => {
        const block:blockInfo = JSON.parse(JSON.parse(event.data))
        const displayVal = blockType[block["bType"]] === "NUMBER" ? block["value"].toString() : blockType[block["bType"]].slice(0, 2)
        modifyBoard(store.board, block["x"], block["y"], displayVal)
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
