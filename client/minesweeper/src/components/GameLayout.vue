<script setup lang='ts'>
import {ref} from 'vue'
const container = ref('main-container')
</script>
<template>
<div :class='container'>
    <slot />
</div>
</template>
<style scoped>
    .main-container {
        flex-grow: 1;
        display: flex;
        flex-direction: row;
        justify-content: center;
        align-items: center;
        column-gap: 2rem;
        height: 100%;
    }
</style>
