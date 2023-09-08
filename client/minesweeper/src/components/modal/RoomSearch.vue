<script setup lang='ts'>
import socket from '@/socket'
import { onMounted } from 'vue';

onMounted(() => {
    socket.send(JSON.stringify({
        name: 'findPublicRoomIds',
        content: JSON.stringify({})
    }))
})
</script>
<template>

</template>
<style>
@import '@/assets/styles/modal.css';
</style>
