<script setup lang='ts'>
import { onMounted, ref, computed, onUnmounted } from 'vue'
import { BOARDSETTING, GAMESTATUS } from '@/config'
import socket from '@/socket'
import {getAlias, setAlias as saveAlias } from '@/docUtils'
import { gameState, roomState, uiState } from '@/store'

const props = defineProps({
    close: {
        type: Function,
        default: () => {}
    }
})

const roomType = ref('public')
const [alias, aliasRef] = [ref(getAlias() || ''), ref<HTMLInputElement>()]
const [passcode, passcodeRef] = [ref(''), ref<HTMLInputElement>()]


const {PLAYER, SIZE, BOMB, TIME_LIMIT} = BOARDSETTING

const player = ref(PLAYER.TWO)
const size = ref(SIZE.MEDIUM)
const bomb = ref(BOMB.NORMAL)
const timeLimit = ref(TIME_LIMIT.NONE)

const createBtn = computed(() => alias.value.length === 0 ? 'btn disabled' : 'btn')

const createRoom = () => {
    gameState.reset();
    saveAlias(alias.value)
    socket.emit('createRoom', {
        alias: alias.value,
        roomType: roomType.value,
        passcode: roomType.value === 'private' ? passcode.value : '',
        player: player.value,
        size: size.value,
        bomb: bomb.value,
        timeLimit: timeLimit.value
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
const setRoomType = (rType: 'private' | 'public') => {
    roomType.value = rType
}
const setPlayer = (num:number) => {
    player.value = num
}
const setSize = (num:number) => {
    size.value = num
}
const setBomb = (num:number) => {
    bomb.value = num
}
const setTimeLimit = (num:number) => {
    timeLimit.value = num
}

const keydownEventHandler = (event: KeyboardEvent) => {
    if (event.key === 'Enter') createRoom()
}

onMounted(() => {
    aliasRef?.value?.addEventListener('keydown', keydownEventHandler)
})

onUnmounted(() => {
    aliasRef?.value?.removeEventListener('keydown', keydownEventHandler)
})
</script>
<template>
    <div class='grid-container'>
        <!-- Alias -->
        <label class='grid-key'>Your Alias</label>
        <div class='grid-value'>
            <input
                ref='aliasRef'
                type='text'
                id='alias'
                class='input autofocus'
                maxlength=12
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
        <!-- Passcode -->
        <template v-if='roomType === `private`'>
            <label class='grid-key'>Passcode</label>
            <div class='grid-value'>
                <input
                    class='input'
                    ref='passcodeRef'
                    type='password'
                    maxlength=4
                    v-model='passcode'
                />
            </div>
        </template>
        <!--# of Players -->
        <label class='grid-key'># of players</label>
        <div class='grid-value btn-group'>
            <button
                v-for='p in PLAYER'
                :key='`PLAYER-${p}`'
                :class='`${player === p ? `btn selected` : `btn`}`'
                @click='() => {setPlayer(p)}'
            >{{ p }}</button>
        </div>
        <!-- Map size -->
        <label class='grid-key'>Map size</label>
        <div class='grid-value btn-group'>
            <button
                v-for='s in SIZE'
                :key='`SIZE-${s}`'
                :class='size === s ? `btn selected` : `btn`'
                @click='() => {setSize(s)}'
            >{{ s }} x {{ s }}</button>
        </div>
        <!-- # of bombs -->
        <label class='grid-key'># of bombs</label>
        <div class='grid-value btn-group'>
            <button
                v-for='b in BOMB'
                :key='`BOMB-${b}`'
                :class='bomb === b ? `btn selected` : `btn`'
                @click='() => {setBomb(b)}'
            >{{ b }}</button>
        </div>
        <!-- Time Limit -->
        <label class='grid-key'>Time Limit</label>
        <div class='grid-value btn-group'>
            <button
                v-for='t in TIME_LIMIT'
                :key='`TIME_LIMIT-${t}`'
                :class='`${timeLimit === t ? `btn selected` : `btn`}`'
                @click='() => {setTimeLimit(t)}'
            >{{ t ? `${t}s` : 'None'}}</button>
        </div>
    </div>
    <div class='modal-row reverse'>
        <button :class='createBtn' @click='createRoom' :disabled='alias.length === 0'>CREATE</button>
        <button class='btn' @click='cancel' :disabled='alias.length === 0'>CANCEL</button>
    </div>
</template>
<style scoped>
@import '@/assets/styles/modal.css';
    
</style>
