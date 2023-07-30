<script setup lang='ts'>
import { onMounted, ref, computed } from 'vue'
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

const { SMALL, MEDIUM, LARGE } = BOARDSETTING.SIZE
const { LITTLE, NORMAL, MANY } = BOARDSETTING.BOMB

const size = ref(MEDIUM)
const bomb = ref(NORMAL)

const createBtn = computed(() => alias.value.length === 0 ? 'btn disabled' : 'btn')

const createRoom = () => {
    saveAlias(alias.value)
    socket.send(JSON.stringify({
        name: 'createRoom',
        content: JSON.stringify({
            alias: alias.value,
            roomType: roomType.value,
            passcode: roomType.value === 'private' ? passcode.value : '',
            size: size.value,
            bomb: bomb.value
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
const setSize = (num:number) => {
    size.value = num
}
const setBomb = (num:number) => {
    bomb.value = num
}

onMounted(() => {
    aliasRef?.value?.addEventListener('keydown', event => {
        if ((event as KeyboardEvent).key === 'Enter') createRoom()
    })
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
                class='autofocus'
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
                    ref='passcodeRef'
                    type='password'
                    maxlength=4
                    v-model='passcode'
                />
            </div>
        </template>
        <!-- Map size -->
        <label class='grid-key'>Map size</label>
        <div class='btn-group'>
            <button
                :class='`${size === SMALL ? `btn selected` : `btn`}`'
                @click='() => {setSize(SMALL)}'
            >{{ SMALL }} x {{ SMALL }}</button>
            <button class='btn'
                :class='`${size === MEDIUM ? `btn selected` : `btn`}`'
                @click='() => {setSize(MEDIUM)}'
            >{{ MEDIUM }} x {{ MEDIUM }}</button>
            <button
                :class='`${size === LARGE ? `btn selected` : `btn`}`'
                @click='() => {setSize(LARGE)}'
            >{{ LARGE }} x {{ LARGE }}</button>
        </div>
        <!-- # of bombs -->
        <label class='grid-key'># of bombs</label>
        <div class='btn-group'>
            <button
                :class='`${bomb === LITTLE ? `btn selected` : `btn`}`'
                @click='() => {setBomb(LITTLE)}'
            >{{ LITTLE }}</button>
            <button
                :class='`${bomb === NORMAL ? `btn selected` : `btn`}`'
                @click='() => {setBomb(NORMAL)}'
            >{{ NORMAL }}</button>
            <button
                :class='`${bomb === LARGE ? `btn selected` : `btn`}`'
                @click='() => {setBomb(MANY)}'
            >{{ MANY }}</button>
        </div>
    </div>
   
    <div class='modal-row reverse'>
        <button :class='createBtn' @click='createRoom' :disabled='alias.length === 0'>CREATE</button>
        <button :class='createBtn' @click='cancel' :disabled='alias.length === 0'>CANCEL</button>
    </div>
</template>
<style scoped>
@import '@/assets/styles/modal.css';
    .grid-container {
        display: grid;
        grid-template-columns: auto auto;
        row-gap: .5rem;
    }
    .grid-key {
        padding-right: 2rem;
        text-align:center;
    }
    .grid-value {
        background: transparent;
        color: white;
        outline-width: 0;
        text-align: center;
        box-sizing: border-box;
    }
    .grid-value input {
        background: transparent;
        border: 1px solid whitesmoke;
        border-radius: .2rem;
        color: white;
        padding: .3rem .3rem .2rem .3rem;
    }
</style>
