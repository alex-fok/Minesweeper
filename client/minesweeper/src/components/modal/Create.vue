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
const [alias, aliasRef] = [ref(getAlias() || ''), ref<HTMLInputElement>()]
const createBtn = computed(() => alias.value.length === 0 ? 'btn disabled' : 'btn')
const setAlias = (event:Event) => {
    alias.value = (event.target as HTMLInputElement).value
}

const createRoom = () => {
    saveAlias(alias.value)
    socket.send(JSON.stringify({
        name: 'createRoom',
        content: JSON.stringify({alias: alias.value})
    }))
    props.close()
}

onMounted(() => {
    aliasRef?.value?.addEventListener('keydown', event => {
        if ((event as KeyboardEvent).key === 'Enter') createRoom()
    })
})
</script>
<template>
    <div class='modal-row'>
        <div class='modal-item'>
            <div class='user-input'>
                <label for='alias'>Your Alias:</label>
                <input
                    ref='aliasRef'
                    type='text'
                    id='alias'
                    class='alias-input autofocus'
                    maxlength=12
                    :value='alias'
                    @input='setAlias'
                />
            </div>
        </div>    
        <div class='modal-item'>
            <button :class='createBtn' @click='createRoom' :disabled='alias.length === 0'>CREATE</button>
        </div>
        <div class='modal-close' @click='close()'>&#10005;</div>
    </div>
</template>
<style scoped>
    @import '@/assets/styles/modal.css';
</style>
