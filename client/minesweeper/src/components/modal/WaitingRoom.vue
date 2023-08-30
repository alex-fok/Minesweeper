<script setup lang='ts'>
import { computed, nextTick, onMounted, onUnmounted, ref } from 'vue'
import socket from '@/socket'
import { gameState } from '@/store'
import Edit from '../icon/edit.vue'

const [alias, setAlias] = [
    ref(gameState.players[gameState.id].alias),
    (str: string) => { alias.value = str }
]

const [isEditVisible, setEditVisibility] = [
    ref(false),
    (bool: boolean) => { isEditVisible.value = bool }
]

const [isEditing, setIsEditing] = [
    ref(false),
    (bool: boolean) => { isEditing.value = bool }
]

const updateReady = () => {
    if (!gameState.isPlayer) return
    socket.send(JSON.stringify({
        name: 'ready',
        content: JSON.stringify({
            isReady: !gameState.players[gameState.id]?.isReady || false
        })
    }))
}
const editStyle = computed(() => ({
    visibility: isEditVisible.value ? 'visible' : 'hidden'
}))

const enableEdit = async() => {
    setIsEditing(true)
    await nextTick()
    const input = document.getElementsByClassName('alias-input')[0] as HTMLInputElement
    input.focus()
}

const endEdit = (event:Event) => {
    if (event.target instanceof Element) {
        if (!event.target.classList.contains('alias-input') &&
        !event.target.classList.contains('self') &&
        isEditing.value) {
            socket.send(JSON.stringify({
                name: 'rename',
                content: JSON.stringify({
                    alias: alias.value
                })
            }))
            isEditing.value = false
        }
    }
}

onMounted(() => {
    const container = document.getElementsByClassName('modal')[0] as HTMLDivElement
    container.addEventListener('click', endEdit)
})
onUnmounted(() => {
    const container = document.getElementsByClassName('modal')[0] as HTMLDivElement
    container?.removeEventListener('click', endEdit)
})
</script>
<template>
    <div class='grid-container'>
        <label class='title'>Players</label>
        <div class='title'>Status</div>
    <template v-for='player in gameState.players'>
        <label
            v-if='gameState.id === player.id && !isEditing'
            class='grid-key'
            :onmouseenter='() => setEditVisibility(true)'
            :onmouseleave='() => setEditVisibility(false)'
            :onclick='() => enableEdit()'
        >
            <span class='self'>
                {{ player.alias }}
                (You)
                <Edit :style='editStyle' fill='white' size='1rem' />
            </span>
            <span :class='player.isOnline ? `online` : `offline`'>({{ player.isOnline ? 'online' : 'offline' }})</span>
        </label>
        <label
            v-else-if='gameState.id === player.id && isEditing'
            class='grid-key'
        >
            <input
                class='alias-input'
                :value='alias'
                @input='event => { setAlias((event.target as HTMLInputElement).value) }'
            />
        </label>
        <label v-else class='grid-key'>
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
    .visible {
        visibility: visible;
    }
    .alias-input {
        background-color: transparent;
        outline: 0;
        color: inherit;
        font-size: 1rem;
        border-style: none none solid;
        border-bottom: 1px white solid;
    }
    .self:hover {
        text-decoration: underline;
        cursor: pointer;
    }
</style>
