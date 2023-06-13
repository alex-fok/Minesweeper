<script setup lang='ts'>
import { ref, onMounted, watch, nextTick } from 'vue';
import { uiState } from '@/store';
import Create from './modal/Create.vue'
import Join from './modal/Join.vue'
import CreateOrJoin from './modal/CreateOrJoin.vue';
import GameEnded from './modal/GameEnded.vue';
import NoContent from './modal/NoContent.vue'

const props = defineProps({
    // create, join, createOrJoin
    content: {
        default: 'createOrJoin'
    },
})

const showingContent = ref(props.content)

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

const setContent = (v:string) => {
    showingContent.value = v
}

onMounted(setFocus)
watch(showingContent, setFocus)
</script>
<template>
    <div class='overlay'></div>
    <div class='modal'>
        <template v-if='showingContent === `create`'>
           <Create :close='close'/> 
        </template>
        <template v-else-if='showingContent === `join`'>
           <Join :close='close'/> 
        </template>
        <template v-else-if='showingContent === `createOrJoin`'>
           <CreateOrJoin :close='close' :setContent='setContent'/>
        </template>
        <template v-else-if='showingContent === `gameEnded`'>
            <GameEnded :close='close'/>
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
