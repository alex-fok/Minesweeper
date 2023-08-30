<script setup lang='ts'>
import { onMounted, watch, nextTick } from 'vue'
import { uiState } from '@/store'
import Create from './modal/Create.vue'
import Join from './modal/Join.vue'
import CreateOrJoin from './modal/CreateOrJoin.vue'
import GameEnded from './modal/GameEnded.vue'
import Invited from './modal/Invited.vue'
import Message from './modal/Message.vue'
import HowToPlay from './modal/HowToPlay.vue'
import Passcode from './modal/Passcode.vue'
import PlayerAlias from './modal/PlayerAlias.vue'
import WaitingRoom from './modal/WaitingRoom.vue'
import NoContent from './modal/NoContent.vue'

const props = defineProps({
    // create, join, createOrJoin
    content: {
        default: 'createOrJoin'
    },
})

// Auto focus when input field is available
const setFocus = async () => {
    await nextTick()
    const elCollection = document.getElementsByClassName('autofocus') as HTMLCollectionOf<HTMLInputElement>
    if (!elCollection.length) return
    elCollection[0].focus()
}

const close = () => {
    uiState.modal.isActive = false 
}

onMounted(setFocus)
watch(props, setFocus)
</script>
<template>
    <div class='overlay'></div>
    <div class='modal'>
        <template v-if='props.content === `create`'>
           <Create :close='close'/> 
        </template>
        <template v-else-if='props.content === `join`'>
           <Join :close='close'/> 
        </template>
        <template v-else-if='props.content === `createOrJoin`'>
           <CreateOrJoin :close='close' />
        </template>
        <template v-else-if='props.content === `gameEnded`'>
            <GameEnded :close='close'/>
        </template>
        <template v-else-if='props.content === `invited`'>
            <Invited :close='close' />
        </template>
        <template v-else-if='props.content === `message`'>
            <Message :close='close' />
        </template>
        <template v-else-if='props.content === `howToPlay`'>
            <HowToPlay :close='close' />
        </template>
        <template v-else-if='props.content === `passcode`'>
            <Passcode :close='close' />
        </template>
        <template v-else-if='props.content === `playerAlias`'>
            <PlayerAlias :close='close' />
        </template>
        <template v-else-if='props.content === `waitingRoom`'>
            <WaitingRoom />
        </template>
        <template v-else>
           <NoContent :close='close'/>
        </template>
    </div>
</template>
<style scoped>
    .overlay {
        position: absolute;
        inset: 0;
        background-color: rgba(207, 207, 207, .4);
        z-index: 1;
    }
    .modal {
        margin: 0;
        position: absolute;
        display: flex;
        flex-direction: column;
        row-gap: 1rem;
        top: 50%;
        left: 50%;
        transform: translate(-50%, -50%);
        z-index: 2;
        padding: 1.5rem 2rem;
        background-color: rgba(52, 52, 52, .95);
        border-radius: .4rem;
        user-select: none;
    }
</style>
