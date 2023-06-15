<script setup lang='ts'>
import Board from './Board.vue'
import Panel from './Panel.vue'
import GameLayout from './GameLayout.vue'
import TopMenu from './TopMenu.vue'
import Modal from './Modal.vue'

import { GAMESTATUS } from '@/config'
import { gameState, uiState } from '@/store'
import { computed } from 'vue'

const { UNDETERMINED, NEW } = GAMESTATUS
const isReady = computed(() => ![UNDETERMINED, NEW].includes(gameState.status))
</script>
<template>
    <Modal
        v-if='uiState.modal.isActive'
        :content='uiState.modal.content !== `` ? uiState.modal.content : undefined'
    />
    <GameLayout>
        <template #header>
            <TopMenu/>
        </template>
        <template #default v-if='isReady'>
            <Board/>
            <Panel/> 
        </template>
    </GameLayout>
</template>
