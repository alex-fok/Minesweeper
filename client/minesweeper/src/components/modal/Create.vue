<script setup lang='ts'>
import { ref } from 'vue'
import socket from '@/socket'

const props = defineProps({
    close: {
        type: Function,
        default: () => {}
    }
})
const alias = ref('')
const createBtn = ref('btn hidden')

const setAlias = (event:Event) => {
    alias.value = (event.target as HTMLInputElement).value
    createBtn.value = alias.value !== '' ? 'btn' : 'btn hidden'
}

const createRoom = () => {
    if (alias.value === '') return
    socket.send(JSON.stringify({
        name: 'createRoom',
        content: JSON.stringify({alias: alias.value})
    }))
    props.close()
}
</script>
<template>
    <div class='modal-item'>
        <div class='info'>
            <label for='alias'>Your Alias:</label>
            <input
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
        <span :class='createBtn' @click='createRoom'>CREATE</span>
    </div>
    <div class='modal-close' @click='close()'>&#10005;</div>
</template>
<style scoped>
    @import '@/assets/modal.css';
</style>
