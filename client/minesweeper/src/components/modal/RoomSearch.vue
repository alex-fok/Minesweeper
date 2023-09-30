<script setup lang='ts'>
import socket from '@/socket'
import { computed, onMounted, ref } from 'vue';
import Reload from '../icon/ReloadIcon.vue';
import { publicState, uiState } from '@/store';
import Avatar from '../icon/AvatarIcon.vue';
const props = defineProps({

    setPrefill: {
        type: Function,
        default: () => {}
    }
})

const [selectedRId, setSelected] = [ref(-1), (num: number) => { selectedRId.value = num }]

const confirmBtn = computed(() => selectedRId.value === -1 ? 'btn disabled' : 'btn')

const findRooms = () => {
    socket.emit('findPublicRoomIds', {})
}

const cancel = () => {
    uiState.modal.displayContent('join')
}

const confirmRoom = () => {
    props.setPrefill('join', { roomId: selectedRId.value })
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
    <template v-if='publicState.rooms.length > 0'>
        <div class='room-list'>
            <div
                v-for='(r, i) in publicState.rooms'
                :key='i'
                :class='selectedRId === r.id ? `room-wrapper selected` : `room-wrapper`'
                @click='() => {
                    if (r.users < r.capacity)
                        setSelected(r.id)
                }'
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
            </div>
        </div>
        <div class='modal-row'>
            <label class=''>Selected:</label>
            <input
                class='room-input'
                type='text'
                :value='selectedRId > -1 ? selectedRId : ``'
                placeholder='Select a room number...'
                disabled />
        </div>
    </template>
    <div v-else class='no-room'>--- No rooms available ---</div>
    <div class='modal-row reverse'>
        <button
            :class='confirmBtn'
            @click='confirmRoom'
            :disabled='selectedRId === -1'
        >CONFIRM</button>
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
    cursor: pointer;
    border: 1px rgba(255, 255, 255, 0) solid;
    padding: .2rem .5rem;
}
.room-wrapper.selected {
    border: 1px white solid;
    border-radius: .3rem;
}
.room-wrapper:hover {
    border: 1px white solid;
    border-radius: .3rem;
}
.room-id {
    width: 5rem;
}
.occupancy {
    display: flex;
    flex-direction: row;
    justify-content: start;
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
