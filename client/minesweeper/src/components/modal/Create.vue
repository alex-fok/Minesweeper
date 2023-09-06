<script setup lang='ts'>
import { onMounted, ref, computed, onUnmounted } from 'vue'
import { BOARDSETTING, GAMESTATUS } from '@/config'
import socket from '@/socket'
import {getAlias, setAlias as saveAlias } from '@/docUtils'
import { gameState, uiState } from '@/store'

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
    socket.send(JSON.stringify({
        name: 'createRoom',
        content: JSON.stringify({
            alias: alias.value,
            roomType: roomType.value,
            passcode: roomType.value === 'private' ? passcode.value : '',
            player: player.value,
            size: size.value,
            bomb: bomb.value,
            timeLimit: timeLimit.value
        })
    }))
    props.close()
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
                :class='`${player === PLAYER.TWO ? `btn selected` : `btn`}`'
                @click='() => {setPlayer(PLAYER.TWO)}'
            >{{ PLAYER.TWO }}</button>
            <button
                :class='`${player === PLAYER.THREE ? `btn selected` : `btn`}`'
                @click='() => {setPlayer(PLAYER.THREE)}'
            >{{ PLAYER.THREE }}</button>
            <button
                :class='`${player === PLAYER.FOUR ? `btn selected` : `btn`}`'
                @click='() => {setPlayer(PLAYER.FOUR)}'
            >{{ PLAYER.FOUR }}</button>
        </div>
        <!-- Map size -->
        <label class='grid-key'>Map size</label>
        <div class='grid-value btn-group'>
            <button
                :class='`${size === SIZE.SMALL ? `btn selected` : `btn`}`'
                @click='() => {setSize(SIZE.SMALL)}'
            >{{ SIZE.SMALL }} x {{ SIZE.SMALL }}</button>
            <button
                :class='`${size === SIZE.MEDIUM ? `btn selected` : `btn`}`'
                @click='() => {setSize(SIZE.MEDIUM)}'
            >{{ SIZE.MEDIUM }} x {{ SIZE.MEDIUM }}</button>
            <button
                :class='`${size === SIZE.LARGE ? `btn selected` : `btn`}`'
                @click='() => {setSize(SIZE.LARGE)}'
            >{{ SIZE.LARGE }} x {{ SIZE.LARGE }}</button>
        </div>
        <!-- # of bombs -->
        <label class='grid-key'># of bombs</label>
        <div class='grid-value btn-group'>
            <button
                :class='`${bomb === BOMB.LITTLE ? `btn selected` : `btn`}`'
                @click='() => {setBomb(BOMB.LITTLE)}'
            >{{ BOMB.LITTLE }}</button>
            <button
                :class='`${bomb === BOMB.NORMAL ? `btn selected` : `btn`}`'
                @click='() => {setBomb(BOMB.NORMAL)}'
            >{{ BOMB.NORMAL }}</button>
            <button
                :class='`${bomb === BOMB.MANY ? `btn selected` : `btn`}`'
                @click='() => {setBomb(BOMB.MANY)}'
            >{{ BOMB.MANY }}</button>
        </div>
        <!-- Time Limit -->
        <label class='grid-key'>Time Limit</label>
        <div class='grid-value btn-group'>
            <button
                :class='`${timeLimit === TIME_LIMIT.NONE ? `btn selected` : `btn`}`'
                @click='() => {setTimeLimit(TIME_LIMIT.NONE)}'
            >None</button>
            <button
                :class='`${timeLimit === TIME_LIMIT.SHORT ? `btn selected` : `btn`}`'
                @click='() => {setTimeLimit(TIME_LIMIT.SHORT)}'
            >{{ TIME_LIMIT.SHORT }}s</button>
            <button
                :class='`${timeLimit === TIME_LIMIT.NORMAL ? `btn selected` : `btn`}`'
                @click='() => {setTimeLimit(TIME_LIMIT.NORMAL)}'
            >{{ TIME_LIMIT.NORMAL }}s</button>
            <button
                :class='`${timeLimit === TIME_LIMIT.LONG ? `btn selected` : `btn`}`'
                @click='() => {setTimeLimit(TIME_LIMIT.LONG)}'
            >{{ TIME_LIMIT.LONG }}s</button>
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
