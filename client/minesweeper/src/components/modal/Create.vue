<script setup lang='ts'>
import { onMounted, ref, computed } from 'vue'
import socket from '@/socket'
import {getAlias, setAlias as saveAlias } from '@/docUtils'

const props = defineProps({
    close: {
        type: Function,
        default: () => {}
    }
})


const roomType = ref('public')
const [alias, aliasRef] = [ref(getAlias() || ''), ref<HTMLInputElement>()]
const [passcode, passcodeRef] = [ref(''), ref<HTMLInputElement>()]
const size = ref(26)
const bomb = ref(100)

const createBtn = computed(() => alias.value.length === 0 ? 'btn disabled' : 'btn')

const createRoom = () => {
    saveAlias(alias.value)
    socket.send(JSON.stringify({
        name: 'createRoom',
        content: JSON.stringify({alias: alias.value})
    }))
    props.close()
}
const setAlias = (event:Event) => {
    alias.value = (event.target as HTMLInputElement).value
}
const setRoomType = (rType: 'private' | 'public') => {
    roomType.value = rType
}
const setPasscode = (event:Event) => {
    passcode.value = (event.target as HTMLInputElement).value
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
                :value='alias'
                @input='setAlias'
            />
        </div>
        <!-- Room Type --> 
        <label class='grid-key'>Room Type</label>
        <div class='grid-value btn-group'>
            <button
                :class='`${roomType === "public" ? "btn selected" : "btn"}`'
                @click='() => {setRoomType(`public`)}'
            >Public</button>
            <button
                :class='`${roomType === "private" ? "btn selected" : "btn"}`'
                @click='() => {setRoomType(`private`)}'
            >Private</button>
        </div>
        <!-- Passcode -->
        <template v-if='roomType === `private`'>
            <label class='grid-key'>Passcode</label>
            <div class='grid-value'>
                <input
                    ref='passcodeRef'
                    type='text'
                    maxlength=4
                    @input='setPasscode'
                />
            </div>
        </template>
        <!-- Map size -->
        <label class='grid-key'>Map size</label>
        <div class='btn-group'>
            <button
                :class='`${size === 16 ? "btn selected" : "btn"}`'
                @click='() => {setSize(16)}'
            >16 x 16</button>
            <button class='btn'
                :class='`${size === 26 ? "btn selected" : "btn"}`'
                @click='() => {setSize(26)}'
            >26 x 26</button>
            <button
                :class='`${size === 36 ? "btn selected" : "btn"}`'
                @click='() => {setSize(36)}'
            >36 x 36</button>
        </div>
        <!-- # of bombs -->
        <label class='grid-key'># of bombs</label>
        <div class='btn-group'>
            <button
                :class='`${bomb === 50 ? "btn selected" : "btn"}`'
                @click='() => {setBomb(50)}'
            >50</button>
            <button
                :class='`${bomb === 100 ? "btn selected" : "btn"}`'
                @click='() => {setBomb(100)}'
            >100</button>
            <button
                :class='`${bomb === 150 ? "btn selected" : "btn"}`'
                @click='() => {setBomb(150)}'
            >150</button>
        </div>
    </div>
   
    <div class='modal-row reverse'>
        <button :class='createBtn' @click='createRoom' :disabled='alias.length === 0'>CREATE</button>
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
    .btn-group {
        display: flex;
        align-items: center;
        justify-content: space-around;
    }
</style>
