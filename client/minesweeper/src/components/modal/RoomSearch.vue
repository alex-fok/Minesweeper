<script setup lang='ts'>
import socket from '@/socket'
import { computed, onMounted, ref } from 'vue';
import Reload from '../icon/ReloadIcon.vue';
import { publicState, uiState } from '@/store';
const props = defineProps({

    setPrefill: {
        type: Function,
        default: () => {}
    }
})

const [selectedIdx, setSelected] = [ref(-1), (num: number) => { selectedIdx.value = num }]

const confirmBtn = computed(() => selectedIdx.value === -1 ? 'btn disabled' : 'btn')

const findRooms = () => {
    socket.emit('findPublicRoomIds', {})
}

const cancel = () => {
    uiState.modal.displayContent('join')
}

const confirmRoom = () => {
    props.setPrefill('join', { roomId: publicState.roomIds[selectedIdx.value] })
    uiState.modal.displayContent('join')
}
onMounted(findRooms)
</script>
<template>
<div class='modal-row'>
    <div class='room-list-title'>Select A Room</div>
    <div class='refresh' title='Reload' @click='findRooms'>
        <Reload color='white' size='1rem'/>
    </div>
</div>
<div class='separator'></div>
<template v-if='publicState.roomIds.length > 0'>
    <div class='room-list'>
        <div
            v-for='(id, i) in publicState.roomIds'
            :key='i'
            :class='selectedIdx === i ? `room-id selected` : `room-id`'
            @click='() => { setSelected(i) }'
        >
            {{ id }}
        </div>
    </div>
    <div class='modal-row'>
        <label class=''>Selected:</label>
        <input
            class='room-input'
            type='text'
            :value='selectedIdx > -1 ? publicState.roomIds[selectedIdx] : ``'
            placeholder='Select a room number...'
            disabled/>
    </div>
</template>
<div v-else class='no-room'>--- No rooms available ---</div>
<div class='modal-row reverse'>
    <button :class='confirmBtn' @click='confirmRoom' :disabled='selectedIdx === -1'>CONFIRM</button>
    <button class='btn' @click='cancel'>CANCEL</button>
</div>
</template>
<style scoped>
@import '@/assets/styles/modal.css';
.room-list-title {
    display: flex;
    justify-content: center;
    align-items: center;
    vertical-align: middle;
    font-size: 1.2rem;
}
.separator {
    border-bottom: 2px solid white;
}
.room-list {
    display: grid;
    grid-template-columns: repeat(4, 5rem);
    row-gap: .2rem;
}
.no-room {
    display: flex;
    align-items: center;
    justify-content: center;
    color: #999999;
    font-style: italic;
}
.room-id {
    display: flex;
    align-items: center;
    column-gap: .2rem;
    cursor: pointer;
}
.room-id.selected {
    font-weight: 700;
}
.room-id::before {
    content: '> ';
    visibility: hidden;
}
.room-id:hover::before {
    visibility: visible;   
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
.room-input {
    background: transparent;
    border: 1px solid #666666;
    border-radius: .2rem;
    width: 12rem;
    box-sizing: border-box;
    color: white;
    padding: .3rem .3rem .2rem .3rem;
}
</style>
