<script setup lang='ts'>
import { ref, onMounted, watch, nextTick } from 'vue';
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

// Auto focus when input field is available
const setFocus = async () => {
    await nextTick()
    const elCollection = document.getElementsByClassName('autofocus') as HTMLCollectionOf<HTMLInputElement>
    if (!elCollection.length) return
    elCollection[0].focus()
}
onMounted(setFocus)
watch(showingContent, setFocus)

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
        content: JSON.stringify({
            id: roomIdInt,
            alias: alias.value
        })
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
               <div class='info'>
                    <label for='alias'>Your Alias:</label>
                    <input
                        type='text'
                        id='alias'
                        class='alias-input autofocus'
                        maxlength=12
                        :value='alias'
                        @input='setAlias' />
                </div>
            </div>
            <div class='modal-item'>
                <span :class='createBtn' @click='createRoom'>CREATE</span>
            </div>
            <div class='modal-close' @click='close()'>&#10005;</div>
        </template>
        <template v-else-if='showingContent === `join`'>
            <div class='modal-item'>
                <div class='info'>
                    <label for='roomId'>Room #</label>
                    <input
                        type='text'
                        id='roomId'
                        class='room-input autofocus'
                        maxlength=4
                        :value='roomId'
                        @input='setRoomId' />
                </div>
                <div class='info'>
                    <label for='alias'>Your Alias:</label>
                    <input
                        type='text'
                        id='alias'
                        class='alias-input'
                        maxlength='12'
                        v-model='alias' />
                </div>
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
        inset: 0;
        background-color: rgba(207, 207, 207, .4);
        z-index: 1;
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
        z-index: 2;
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
        column-gap: 1rem;
    }
    .info {
        display: flex;
        flex-direction: column;
        row-gap: .5rem;
    }

    .info:not(:first-child) {
        border-left: 1px solid #9F9F9F;
        padding-left: 1rem;
    }

    .info label {
        margin-right: .5rem;
    }
    .info input {
        background: transparent;
        border: 0;
        color: white;
        outline-width: 0;
        text-align: center;
        box-sizing: border-box;
    }
    .info input:not(:focus) {
        border-bottom: 1px solid #9F9F9F;
        margin-bottom: -.8px;
    }
    .alias-input {
        width: 8rem;
    }
    .room-input {
        width: 4rem;
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
