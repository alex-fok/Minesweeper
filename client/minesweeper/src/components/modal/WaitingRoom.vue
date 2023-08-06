<script setup lang='ts'>
import socket from '@/socket';
import { gameState } from '@/store';

const updateReady = () => {
    if (!gameState.isPlayer) return
    socket.send(JSON.stringify({
        name: 'ready',
        content: JSON.stringify({
            isReady: !gameState.players[gameState.id]?.isReady || false
        })
    }))
}
</script>
<template>
    <div class='grid-container'>
        <label class='title'>Players</label>
        <div class='title'>Status</div>
    <template v-for='player in gameState.players'>
        <label class='grid-key'>
            {{ player.alias }}
            <span :class='player.isOnline ? `online` : `offline`'>({{ player.isOnline ? 'online' : 'offline' }})</span>
        </label>
        <div class='grid-value'>{{ player.isReady ? 'Ready' : 'Not Ready' }}</div>
    </template>
    </div>
    <div v-if= 'gameState.isPlayer' class='modal-row reverse'>
        <button class='btn' @click='updateReady'>
            {{ gameState.players[gameState.id]?.isReady ? 'Unready' : 'I\'m Ready!' }}
        </button>
    </div>
</template>
<style scoped>
    @import '@/assets/styles/modal.css';
    .online {
        color:aquamarine;
    }
    .offline {
        color:crimson;
    }
</style>
