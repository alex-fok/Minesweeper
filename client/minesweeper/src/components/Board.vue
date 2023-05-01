<script setup lang='ts'>
    import { io } from 'socket.io-client'
    const socket = io('http://localhost:8080')

</script>
<template>
    <div class='container'>
       <div
          v-for='(block, i) in board'
          @click='reveal(i)'
          :key='i'
        >
        </div>
    </div>
</template>
<script lang='ts'>
    export default {
        name: 'Board',
        data() {
            return {
                board: new Array(256)
            }
        },
        methods: {
            reveal(num: any) {
                console.log(num)
            }
        }
    }
</script>
<style scoped>
    .container {
        width:96vh;
        height:96vh;
        column-gap: 1px;
        row-gap:1px;
        display: grid;
        grid-template-columns:auto auto auto auto auto auto auto auto auto auto auto auto auto auto auto auto;
    }
    .container > div {
        background-color:lightblue ;
        
        width: 6vh;
        height: 6vh;
    }
    
</style>
