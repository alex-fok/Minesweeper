<script setup lang='ts'>
import { onMounted, nextTick, ref, watch, onUnmounted } from 'vue'
import { uiState } from '@/store'
import Create from './modal/Create.vue'
import Join from './modal/Join.vue'
import RoomSearch from './modal/RoomSearch.vue'
import CreateOrJoin from './modal/CreateOrJoin.vue'
import GameEnded from './modal/GameEnded.vue'
import Invited from './modal/Invited.vue'
import Message from './modal/Message.vue'
import HowToPlay from './modal/HowToPlay.vue'
import Passcode from './modal/Passcode.vue'
import PlayerAlias from './modal/PlayerAlias.vue'
import WaitingRoom from './modal/WaitingRoom.vue'
import NoContent from './modal/NoContent.vue'

type EventListener = {
    event: string,
    handler: (ev : Event) => void
}

type JoinPrefillContent = {
    alias: string,
    roomType: 'private' | 'public',
    roomId: string,
    passcode: string
}

type PrefillType = '' | 'join'
type PrefillContent = {} | Partial<JoinPrefillContent>

type Prefill = {
    type: PrefillType,
    content: PrefillContent
}
const props = defineProps({
    // create, join, createOrJoin
    content: {
        type: String,
        default: 'createOrJoin'
    }
})
const modalRef = ref<HTMLDivElement>()
const prefill = ref<Prefill>({
    type: '',
    content: {}
})

const setPrefill = (type: PrefillType, content: PrefillContent) => {
    console.log('type:', type)
    console.log('content:', content)
    prefill.value.type = type
    prefill.value.content = { ...prefill.value.content, ...content }
}

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
const eventListeners: EventListener[] = []

const addModalEventListener = (eventListener: EventListener) => {
    const { event, handler } = eventListener
    modalRef.value?.addEventListener(event, handler)
    eventListeners.push(eventListener)
}

onMounted(setFocus)

onUnmounted(() => {
    eventListeners.forEach(eventListener => {
        const { event, handler } = eventListener
        modalRef.value?.removeEventListener(event, handler)
    })
})
watch(props, setFocus)
</script>
<template>
    <div class='overlay'></div>
    <div ref='modalRef' class='modal'>
        <template v-if='props.content === `create`'>
           <Create :close='close'/> 
        </template>
        <template v-else-if='props.content === `join`'>
           <Join
                :prefill='prefill.type === `join` ? prefill.content : {}'
                :setPrefill='setPrefill'
                :close='close'
            />
        </template>
        <template v-else-if='props.content === `roomSearch`'>
            <RoomSearch :setPrefill='setPrefill' />
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
            <WaitingRoom :addModalEventListener='addModalEventListener'/>
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
