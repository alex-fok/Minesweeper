<script setup lang='ts'>
    import { ref } from 'vue'
    defineProps({
        turnCount: Number 
    })
    const panel = ref('side-container')
</script>
<template>
    <div :class='panel'>
        <h2>Turn {{ turnCount }}</h2>
        <br/>
        
    </div>
</template>
<style scoped>
    .side-container {
        margin-top: 0%;
        color:white;
        padding: 0em 2em;
        box-sizing:border-box;
        background-color: #343434;
        width: 15rem;
        height: 100%;
    }
</style>
