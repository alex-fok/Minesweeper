<script setup lang='ts'>
import { ref } from 'vue';
import { socket } from '@/socket';

const props = defineProps({
    close: {
        type: Function,
        default: () => {}
    }
})

const alias = ref('')
const roomId = ref('')
const joinBtn = ref('btn hidden')

const setRoomId = (event:Event) => {
    roomId.value = (event.target as HTMLInputElement).value
    joinBtn.value = roomId.value.length ? 'btn' : 'btn hidden'
}

const joinRoom = () => {
    const roomIdInt = parseInt(roomId.value, 10)
    if (Number.isNaN(roomIdInt)) return
    socket.send(JSON.stringify({
        name: 'joinRoom',
        content: JSON.stringify({
            id: roomIdInt,
            alias: alias.value
        })
    }))
    props.close()
}
</script>
<template>
    <div class='modal-item'>
        <div class='info'>
            <label for='roomId'>Room #</label>
            <input
                type='text'
                id='roomId'
                class='room-input autofocus'
                maxlength=4
                :value='roomId'
                @input='setRoomId'
            />
        </div>
        <div class='info'>
            <label for='alias'>Your Alias:</label>
            <input
                type='text'
                id='alias'
                class='alias-input'
                maxlength='12'
                v-model='alias'
            />
        </div>
    </div>
    <div class='modal-item'>
        <span :class='joinBtn' @click='joinRoom'>JOIN</span>
    </div>
    <div class='modal-close' @click='close()'>&#10005;</div>
</template>
<style scoped>
    @import '@/assets/modal.css';
</style>
