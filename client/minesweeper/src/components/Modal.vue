<script setup lang='ts'>
import { ref, onUpdated, onMounted } from 'vue';
import { socket } from '@/socket';
import { uiState } from '@/store';

const props = defineProps({
    // create, join, createOrJoin
    content: {
        default: 'createOrJoin'
    },
})

const alias = ref('')
const roomId = ref('')
const createBtn = ref('btn hidden')
const joinBtn = ref('btn hidden')
const showingContent = ref(props.content)

// Auto Focus when input field is available
const setFocus = () => {
    const elCollection = document.getElementsByClassName('autofocus') as HTMLCollectionOf<HTMLInputElement>
    if (!elCollection.length) return
    elCollection[0].focus()
}
onUpdated(() => setFocus())
onMounted(() => setFocus())

const close = () => {
    uiState.active = true
}

const createRoom = () => {
    if (alias.value === '') return
    socket.send(JSON.stringify({
        name: 'createRoom',
        content: JSON.stringify({alias: alias.value})
    }))
    close()
}

const joinRoom = () => {
    const roomIdInt = parseInt(roomId.value, 10)
    if (Number.isNaN(roomIdInt)) return
    socket.send(JSON.stringify({
        name: 'joinRoom',
        content: JSON.stringify({id: roomIdInt})
    }))
    close()
}

const setAlias = (event:Event) => {
    alias.value = (event.target as HTMLInputElement).value
    createBtn.value = alias.value !== '' ? 'btn' : 'btn hidden'
}

const setRoomId = (event:Event) => {
    roomId.value = (event.target as HTMLInputElement).value
    joinBtn.value = roomId.value.length ? 'btn' : 'btn hidden'
}

const setContent = (v:string) => {
    showingContent.value = v
}
</script>
<template>
    <div class='overlay'></div>
    <div class='modal'>
        <template v-if='showingContent === `create`'>
            <div class='modal-item'>
               <span>
                    <label for='alias'>Your Alias:</label>
                    <input
                        type='text'
                        id='alias'
                        class='alias-input autofocus'
                        maxlength=12
                        :value='alias'
                        @input='setAlias' />
                </span>
            </div>
            <div class='modal-item'>
                <span :class='createBtn' @click='createRoom'>CREATE</span>
            </div>
            <div class='modal-close' @click='close()'>&#10005;</div>
        </template>
        <template v-else-if='showingContent === `join`'>
            <div class='modal-item'>
                <span>
                    <label for='roomId'>Room #</label>
                    <input
                        type='text'
                        id='roomId'
                        class='room-input autofocus'
                        maxlength=4
                        :value='roomId'
                        @input='setRoomId' />
                </span>
            </div>
            <div class='modal-item'>
                <span :class='joinBtn' @click='joinRoom'>JOIN</span>
            </div>
            <div class='modal-close' @click='close'>&#10005;</div>
        </template>
        <template v-else-if='showingContent == `createOrJoin`'>
            <div class='modal-item'>
                <span class='btn' @click='setContent(`create`)'>CREATE ROOM</span>
            </div>
            <div class='modal-item'>
                OR
            </div>
            <div class='modal-item'>
                <span class='btn' @click='setContent(`join`)'>JOIN ROOM</span>
            </div>
            <div class='modal-close' @click='close'>&#10005;</div>
        </template>
        <template v-else>
            No Content
            <div class='modal-close' @click='close'>&#10005;</div>
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
        align-items: stretch;
        column-gap: 1.5rem;
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
    }
    .modal-item {
        flex-grow: 1;
        display: flex;
        margin: auto;
    }
    .modal-item input {
        background: transparent;
        border: 0;
        color: white;
        margin: 0 .5rem;
        outline-width: 0;
    }
    .alias-input {
        width: 8rem;
    }
    .room-input {
        width: 3rem;
    }
    .btn {
        font-size: .8rem;
        cursor: pointer;
    }
    .btn:hover {
        border-bottom: 1px solid white;
        margin-bottom: -1px;
    }
    .hidden {
        visibility: hidden;
    }
</style>
