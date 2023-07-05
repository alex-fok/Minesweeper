<script setup lang='ts'>
import { ref, computed, onMounted } from 'vue';
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

const setRoomId = (event:Event) => {
    roomId.value = (event.target as HTMLInputElement).value
}

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

onMounted(() => {
    [aliasRef, roomIdRef].forEach(ref => {
        ref?.value?.addEventListener('keydown', event => {
            if ((event as KeyboardEvent).key === 'Enter') joinRoom()
        })
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
                    :value='roomId'
                    @input='setRoomId'
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
