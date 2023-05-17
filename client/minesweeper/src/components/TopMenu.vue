<script setup lang='ts'>
import { ref } from 'vue'
const container = ref('header-container')
defineProps({
    roomId: Number
})
</script>
<template>
    <div :class='container '>
        <div><span>Room #{{ roomId && roomId >= 0 ? roomId : ' (Loading...)' }}</span></div>
    </div>
</template>
<style scoped>
    .header-container {
        color: white;
        height: 3rem;
        border-bottom: .1rem solid #9F9F9F;
        padding: .5rem 2rem;
        position: relative;
    }
    .header-container > div {
        position: absolute;
        top: 50%;
        height: 2rem;
        font-size: 1.2rem;
        margin-top: -1rem;
        margin-bottom: 0rem;
    }
    .header-container > div > span {
        line-height: 2rem;
        vertical-align: middle;
    }
</style>
