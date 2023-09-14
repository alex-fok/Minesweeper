<script setup lang='ts'>
import Board from './components/GameBoard.vue'
import GameInfo from './components/GameInfo.vue'
import GameLayout from './components/GameLayout.vue'
import TopMenu from './components/TopMenu.vue'
import ModalView from './components/ModalView.vue'

import { GAMESTATUS } from '@/config'
import { roomState, uiState } from '@/store'
import { computed } from 'vue'

const { UNDETERMINED, NEW, INVITED } = GAMESTATUS
const isReady = computed(() => ![UNDETERMINED, NEW, INVITED].includes(roomState.status))
</script>
<template>
    <ModalView
        v-if='uiState.modal.isActive'
        :content='uiState.modal.content !== `` ? uiState.modal.content : undefined'
    />
    <GameLayout>
        <template #header>
            <TopMenu/>
        </template>
        <template #default v-if='isReady'>
            <Board/>
            <GameInfo/> 
        </template>
    </GameLayout>
</template>
