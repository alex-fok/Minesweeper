<script setup lang='ts'>
import Board from './Board.vue'
import Panel from './Panel.vue'
import GameLayout from './GameLayout.vue'
import TopMenu from './TopMenu.vue'
import Modal from './Modal.vue'

import { GAMESTATUS } from '@/config'
import { gameState, uiState } from '@/store'

const { NEW } = GAMESTATUS
</script>
<template>
    <Modal v-if='uiState.modal.isActive' :content='uiState.modal.content !== `` ? uiState.modal.content : undefined'/>
    <GameLayout>
        <template #header>
            <TopMenu/>
        </template>
        <template #default v-if='gameState.status !== NEW'>
            <Board/>
            <Panel/> 
        </template>
    </GameLayout>
</template>
