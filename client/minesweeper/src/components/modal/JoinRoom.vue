<script setup lang='ts'>
import { ref, computed, onMounted } from 'vue'
import { GAMESTATUS } from '@/config'
import { roomState, uiState, publicState } from '@/store'
import socket from '@/socket';
import { getAlias, setAlias as saveAlias } from '@/docUtils'
import Reload from '../icon/ReloadIcon.vue'
import Avatar from '../icon/AvatarIcon.vue'

const props = defineProps({
    close: {
        type: Function,
        default: () => {}
    }
})

const alias = ref(getAlias() || '')
const roomType = ref('public')
const selected = ref(-1)
const passcode = ref('')

const hoveringRoomId = ref(-1)
const selectedInput = ref<HTMLInputElement>();

const serarchOpenRooms = () => {
    socket.emit('findPublicRoomIds', {})
}
const isJoinable = computed(() =>

    (alias.value.length === 0 || selected.value === -1)
)

const joinRoom = () => {
    if (alias.value !== '') saveAlias(alias.value)

    socket.emit('joinRoom', {
        alias: alias.value,
        roomType: roomType.value,
        id: selected.value,
        passcode: roomType.value === 'private' ? passcode.value : ''
    })
    props.close()
}

const cancel = () => {
    if (roomState.status === GAMESTATUS.NEW) {
        uiState.modal.displayContent('createOrJoin')
    } else {
        props.close()
    }
}

const isHighlighted = (id: number, index: number) => {
    const isSelected = selected.value === id
    const isHovering = hoveringRoomId.value === id
    const isFull = publicState.rooms[index].users === publicState.rooms[index].capacity
   
    return isSelected || (isHovering && !isFull)
}

const setSelected = (id: number) => {
    if (id > 9999) {
        if (selectedInput.value)
            selectedInput.value.value = selected.value.toString()
        return
    }

    const isNaN = Number.isNaN(id)
    
    selected.value = isNaN ? -1 : id;
    
    if (selectedInput.value)
        selectedInput.value.value = isNaN ? '' : id.toString()
}

const setRoomType = (rType: 'private' | 'public') => {
    roomType.value = rType
}

const setHoveringRoom = (id: number) => {
    hoveringRoomId.value = id
}

onMounted(serarchOpenRooms)
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
            <input
                ref='selectedInput'
                class='input autofocus'
                type='number'
                @input='event => { setSelected(parseInt((event.target as HTMLInputElement).value)) }'
            />
        </div>
        <template v-if='roomType === `private`'>
            <label class='grid-key'>Passcode</label>
            <div class='grid-value'>
                <input
                    class='input'
                    type='password'
                    v-model='passcode'
                    maxlength='4'
                />
            </div>
        </template>
    </div>
    <template v-if='roomType === `public`'>
            <div class='modal-row'>
                <div class='room-list-title'>Select A Room</div>
                <div class='refresh' title='Reload' @click='serarchOpenRooms'>
                    <Reload color='white' size='1rem'/>
                </div>
            </div>
            <div class='separator'></div>
            <template v-if='publicState.rooms.length > 0'>
            <div class='room-list'>
                <div
                    v-for='(r, i) in publicState.rooms'
                    :key='i'
                    :class='isHighlighted(r.id, i) ? `room-wrapper highlight` : `room-wrapper`'
                    @mouseover='() => { setHoveringRoom(r.id) }'
                    @mouseout='() => { setHoveringRoom(-1) }'
                    @click='() => { if (r.users < r.capacity) setSelected(r.id) }'
                >
                    <span class='room-id'>{{ r.id }}</span>
                    <span class='occupancy'>
                        <template v-for='c in r.capacity' :key='`capacity-${c}`'>
                            <Avatar
                                size='1rem'
                                :color='c <= r.users ? `white` : `rgba(256, 256, 256, .3)`'
                            />
                        </template>
                    </span>
                    <span v-if='r.users === r.capacity' class='room-full'>(FULL)</span>
                </div>
            </div>
            </template>
            <div v-else class='no-room'>--- No rooms available ---</div>
        </template>
    <div class='modal-row reverse'>
        <button :class='isJoinable ? `btn disabled` : `btn`' @click='joinRoom' :disabled='isJoinable'>JOIN</button>
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
    .room-list-title {
        display: flex;
        justify-content: center;
        align-items: center;
        vertical-align: middle;
        font-size: 1.2rem;
    }
    .refresh {
        color: white;
        background-color: transparent;
        border: 0;
        padding: 0;
        align-self: center;
        cursor: pointer;
        vertical-align: middle;
    }
    .separator {
        border-bottom: 2px solid white;
    }
    .room-list {
        display: flex;
        flex-direction: column;
    }
    .no-room {
        display: flex;
        align-items: center;
        justify-content: center;
        color: #999999;
        font-style: italic;
    }
    .room-wrapper {
        display: flex;
        align-items: center;
        border: 1px rgba(255, 255, 255, 0) solid;
        padding: .2rem .5rem;
    }
    .highlight {
        border: 1px white solid;
        border-radius: .3rem;
    }
    .room-wrapper.highlight {
        border: 1px white solid;
        border-radius: .3rem;
        cursor: pointer;
    }
    .room-id {
        width: 5rem;
    }
    .occupancy {
        display: flex;
        flex-direction: row;
        justify-content: start;
    }
    .room-full {
        color:crimson;
        margin-left: 1.5rem;
    }

    input::-webkit-outer-spin-button,
    input::-webkit-inner-spin-button {
        -webkit-appearance: none;
        margin: 0;
    }

    /* Firefox */
    input[type=number] {
        -moz-appearance: textfield;
    }
</style>
