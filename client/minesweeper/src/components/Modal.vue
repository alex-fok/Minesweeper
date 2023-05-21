<script setup lang='ts'>
import { socket } from '@/socket';
import { ref } from 'vue';

const roomId = ref('')
const joinBtn = ref('join-btn hidden')
const test = ref(true)

const props = defineProps({
    content: {
        default: 'joinRoom'
    } 
})
const joinRoom = () => {
    const roomIdInt = parseInt(roomId.value, 10)
    if (Number.isNaN(roomIdInt)) return
    socket.send(JSON.stringify({
        name: 'joinRoom',
        content: JSON.stringify({id: roomIdInt})
    }))
}

const setRoomId = (event:Event) => {
    roomId.value = (event.target as HTMLInputElement).value
    joinBtn.value = roomId.value.length ? 'join-btn' : 'join-btn hidden'
}

</script>
<template>
    <div class='overlay'></div>
    <div class='modal'>
        <template v-if='content === `joinRoom`'>
            <div class='modal-item'>
                <span>
                    <label for='roomId'>Room #</label>
                    <input
                        type='text'
                        id='roomId'
                        maxlength='4'
                        :value='roomId'
                        @input='setRoomId'
                        autofocus />
                </span>
            </div>
            <div class='modal-item'>
                <span :class='joinBtn' @click='joinRoom'>JOIN</span>
            </div>
            <div class='modal-close'>&#10005;</div>
        </template>
        <template v-else-if='content === `createOrJoin`'>
            Create Or Join 
            <div class='modal-close'>&#10005;</div>
        </template>
        <template v-else>
            No Content
        </template>
    </div>
</template>
<style scoped>
    .overlay {
        position: absolute;
        width: 100vw;
        height: 100vh;
        background-color: rgba(207, 207, 207, .4);
    }
    .modal {
        margin: 0;
        position: absolute;
        display: flex;
        flex-direction: row;
        top: 50%;
        left: 50%;
        transform: translate(-50%, -50%);
        padding: 1.5rem 2rem;
        background-color: rgba(52, 52, 52, .95);
        border-radius: .4rem;
        user-select: none;
    }
    .modal-close {
        align-self: center;
        cursor: pointer;
        vertical-align: middle;
        margin-left: 0.7rem;
    }
    .modal-item {
        flex-grow: 1;
        display: flex;
        align-items: center;
        margin: auto;
    }
    .modal-item input {
        background: transparent;
        border: 0;
        color: white;
        margin: 0 .5rem;
        width: 4rem;
        box-sizing: content-box;
        outline-width: 0;
    }
    
    .join-btn {
        font-size: .8rem;
        cursor: pointer;
    }

    .hidden {
        visibility: hidden;
    }
</style>
