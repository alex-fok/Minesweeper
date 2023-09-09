<script setup lang='ts'>
import { ref, computed, onMounted, onUnmounted } from 'vue';
import { GAMESTATUS } from '@/config';
import { gameState, uiState } from '@/store';
import socket from '@/socket';
import { getAlias, setAlias as saveAlias } from '@/docUtils'
import Search from '../icon/search.vue';

const props = defineProps({
    prefill: {
        type: Object,
        default: {}
    },
    setPrefill: {
        type: Function,
        default: () => {}
    },
    close: {
        type: Function,
        default: () => {}
    }
})

const { prefill } = props

const alias = ref(prefill.alias || getAlias() || '')
const aliasRef = ref<HTMLInputElement>()

const roomType = ref(prefill.roomType || 'public')

const roomId = ref(prefill.roomId || '')
const roomIdRef = ref<HTMLInputElement>()

const passcode = ref(prefill.passcode || '')
const passcodeRef = ref<HTMLInputElement>()

const isFindRoomFocus = ref(false)
const setIsFindRoomFocus = (bool: boolean) => {
    isFindRoomFocus.value = bool
}

const joinBtn = computed(() => 
    (alias.value.length === 0 || roomId.value.length === 0) ? 'btn disabled' : 'btn'
)

const joinRoom = () => {
    const roomIdInt = parseInt(roomId.value, 10)
    if (Number.isNaN(roomIdInt)) return

    if (alias.value !== '') saveAlias(alias.value)

    socket.send(JSON.stringify({
        name: 'joinRoom',
        content: JSON.stringify({
            alias: alias.value,
            roomType: roomType.value,
            id: roomIdInt,
            passcode: roomType.value === 'private' ? passcode.value : ''
        })
    }))
    props.close()
}

const search = () => {
    console.log('search')
    props.setPrefill('join', {
        alias: alias.value,
        roomType: roomType.value,
        roomId: roomId.value,
        passcode: passcode.value
    })
    uiState.modal.displayContent('roomSearch')
}
const cancel = () => {
    if (gameState.status === GAMESTATUS.NEW) {
        uiState.modal.displayContent('createOrJoin')
    } else {
        props.close()
    }
}
const setRoomType = (rType: 'private' | 'public') => {
    roomType.value = rType
}

const keydownEventHandler = (event: KeyboardEvent) => {
    if (event.key === 'Enter') joinRoom()
}

onMounted(() => {
    [aliasRef, roomIdRef, passcodeRef].forEach(ref => {
        ref?.value?.addEventListener('keydown', keydownEventHandler)
    })
})
onUnmounted(() => {
    [aliasRef, roomIdRef, passcodeRef].forEach(ref => {
        ref?.value?.removeEventListener('keydown', keydownEventHandler)
    })
})
</script>
<template>
    <div class='grid-container'>
        <!-- Alias -->
        <label class='grid-key'>Your Alias</label>
        <div class='grid-value'>
            <input
                class='input autofocus'
                ref='aliasRef'
                type='text'
                id='alias'
                v-model='alias'
            />
        </div>
        <!-- Room Type --> 
        <label class='grid-key'>Room Type</label>
        <div class='grid-value btn-group'>
            <button
                :class='`${roomType === `public` ? `btn selected` : `btn`}`'
                @click='() => {setRoomType(`public`)}'
            >Public</button>
            <button
                :class='`${roomType === `private` ? `btn selected` : `btn`}`'
                @click='() => {setRoomType(`private`)}'
            >Private</button>
        </div>
        <label class='grid-key'>Room #</label>
        <div class='grid-value'>
            <div
                v-if='roomType === `public`'
                class='search-container'
            >
                <input
                    class='input'
                    v-model='roomId'
                    maxlength='4'
                    @focus='() => { setIsFindRoomFocus(true) }'
                    @blur='() => { setIsFindRoomFocus(false) }'
                />
                <!-- Dropdown -->
                <span v-if='isFindRoomFocus' class='search-wrapper'>
                    <Search color='white' size='2vh'/>
                    <div
                        class='search-dropdown'
                        @mousedown='search'
                    >Find Public Room...</div>
                </span>
            </div>
            <input
                v-else
                class='input autofocus'
                ref='roomIdRef'
                type='text'
                v-model='roomId'
                maxlength='4'
            />
        </div>
        <template v-if='roomType === `private`'>
            <label class='grid-key'>Passcode</label>
            <div class='grid-value'>
                <input
                    class='input'
                    ref='passcodeRef'
                    type='password'
                    v-model='passcode'
                    maxlength='4'
                />
            </div>
        </template>
    </div>
    <div class='modal-row reverse'>
        <button :class='joinBtn' @click='joinRoom' :disabled='alias.length === 0'>JOIN</button>
        <button class='btn' @click='cancel'>CANCEL</button>
    </div>
</template>
<style scoped>
    @import '@/assets/styles/modal.css';
    .search-container {
        position:relative;
        width: 12rem;
    }
    .search-dropdown {
        display: flex;
        border: 0;
        color: inherit;
    }
    .search-wrapper {
        position: absolute;
        display: flex;
        flex-direction: row;
        width: 12rem;
        column-gap: 1rem;
        background-color: #1C1C1C;
        box-sizing:border-box;
        padding: .5rem .5rem;
        cursor: pointer;
        z-index: 1;
    }
    .search-wrapper:hover {
        background-color: #666666;
    }
    .search-container .input {
        outline: none;
    }
</style>
