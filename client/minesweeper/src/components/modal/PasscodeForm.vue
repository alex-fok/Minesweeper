<script setup lang='ts'>
import socket from '@/socket';
import { onMounted, onUnmounted, ref } from 'vue';
const props = defineProps({
    close: {
        type: Function,
        default: () => {}
    }
})
const [passcode, passcodeRef] = [ref(''), ref<HTMLInputElement>()]

const sendPass = () => {
    socket.emit('passcode', { passcode: passcode.value })
    props.close()
}
const keyDownEventHandler = (event: KeyboardEvent) => {
    if (event.key === 'Enter') {
        sendPass()
        props.close()
    }   
}
onMounted(() => {
    passcodeRef.value?.addEventListener('keydown', keyDownEventHandler)
})
onUnmounted(() => {
    passcodeRef?.value?.removeEventListener('keydown', keyDownEventHandler)
})
</script>
<template>
    <div class='modal-row'>
        <div class='modal-item'>
            <div class='user-input'>
                <label for='passcode'>Passcode</label>
                <input
                    ref='passcodeRef'
                    type='password'
                    id='passcode'
                    class='autofocus'
                    maxlength=4
                    v-model='passcode'
                />
            </div>
        </div>
        <div class='modal-item'>
            <button class='btn' @click='sendPass'>CONFIRM</button>
        </div>
        <div class='modal-close' @click='close()'>&#10005;</div>
    </div>
</template>
<style>
    @import '@/assets/styles/modal.css';
</style>
