<script setup lang='ts'>
import { onMounted, ref } from 'vue';
import { gameState } from '@/store';
import socket from '@/socket';

const props = defineProps({
    close: {
        type: Function,
        default: () => {}
    }
})
const cancel = () => {
    props.close()
}
const alias = ref(gameState.players[gameState.id]?.alias || '')
const aliasRef = ref<HTMLInputElement>()

const setAlias = (event: Event) => {
    alias.value = (event.target as HTMLInputElement).value
}

const rename = () => {
    socket.send(JSON.stringify({
        name: 'rename',
        content: JSON.stringify({
            alias: alias.value
        })
    }))
    props.close()
}

onMounted(() => {
    aliasRef?.value?.addEventListener('keydown', event => {
        if ((event as KeyboardEvent).key === 'Enter') rename()
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
    </div>
    <div class='modal-row reverse'>
        <button class='btn' @click='rename' :disabled='alias.length === 0'>RENAME</button>
        <button class='btn' @click='cancel'>CANCEL</button>
    </div>
</template>
<style>
@import '@/assets/styles/modal.css';
</style>
