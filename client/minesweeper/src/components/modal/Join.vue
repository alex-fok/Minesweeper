<script setup lang='ts'>
import { ref, computed, onMounted, onUnmounted } from 'vue';
import socket from '@/socket';
import { getAlias, setAlias as saveAlias } from '@/docUtils'

const props = defineProps({
    close: {
        type: Function,
        default: () => {}
    }
})

const [alias, aliasRef] = [ref(getAlias() || ''), ref<HTMLInputElement>()]
const [roomId, roomIdRef] = [ref(''), ref<HTMLInputElement>()]

const joinBtn = computed(() => (alias.value.length === 0 || roomId.value.length === 0) ? 'btn disabled' : 'btn')

const joinRoom = () => {
    const roomIdInt = parseInt(roomId.value, 10)
    if (Number.isNaN(roomIdInt)) return

    if (alias.value !== '') saveAlias(alias.value)

    socket.send(JSON.stringify({
        name: 'joinRoom',
        content: JSON.stringify({
            id: roomIdInt,
            alias: alias.value
        })
    }))
    props.close()
}

const keyDownEventHandler = (event: KeyboardEvent) => {
    if (event.key === 'Enter') joinRoom()
}
onMounted(() => {
    [aliasRef, roomIdRef].forEach(ref => {
        ref?.value?.addEventListener('keydown', keyDownEventHandler)
    })
})
onUnmounted(() => {
    [aliasRef, roomIdRef].forEach(ref => {
        ref?.value?.removeEventListener('keydown', keyDownEventHandler)
    })
})
</script>
<template>
    <div class='modal-row'>
        <div class='modal-item'>
            <div class='user-input'>
                <label for='roomId'>Room #</label>
                <input
                    ref='roomIdRef'
                    type='text'
                    id='roomId'
                    class='room-input autofocus'
                    maxlength=4
                    v-model='roomId'
                />
            </div>
            <div class='user-input'>
                <label for='alias'>Your Alias:</label>
                <input
                    ref='aliasRef'
                    type='text'
                    id='alias'
                    class='alias-input'
                    maxlength='12'
                    v-model='alias'
                />
            </div>
        </div>
        <div class='modal-item'>
            <button :class='joinBtn' @click='joinRoom'>JOIN</button>
        </div>
        <div class='modal-close' @click='close()'>&#10005;</div>
    </div>
</template>
<style scoped>
    @import '@/assets/styles/modal.css';
</style>
