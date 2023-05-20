<script setup lang='ts'>
import { ref } from 'vue'
const container = ref('header-container')
const room = ref('room')
const option = ref('option')
const button = ref('button')
defineProps({
    roomId: Number
})
</script>
<template>
    <div :class='container'>
        <div :class='option'><div :class='room'>Room #{{ roomId && roomId >= 0 ? roomId : ' (Loading...)' }}</div></div>
        <div :class='option'><span :class='button'>Create Game</span></div>
        <div :class='option'><span :class='button'>Join Game</span></div>
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
    .option {
        position: relative;
        height: 4rem;
        line-height: 4rem;
        vertical-align: middle;
        padding: 0rem 1rem;
    }
    .room {
        font-size: 1.2rem;
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
