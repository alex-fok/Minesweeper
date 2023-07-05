<script setup lang='ts'>
import { uiState } from '@/store';

defineProps({
    close: {
        type: Function,
        default: () => {}
    }
})
</script>
<template>
    <div class='modal-row'>
        <div class='modal-item'>
            <button class='btn' @click='uiState.modal.displayContent(`create`)'>CREATE ROOM</button>
        </div>
        <div class='modal-item'>
            <button class='btn' @click='uiState.modal.displayContent(`join`)'>JOIN ROOM</button>
        </div>
        <div class='modal-item'>
            <button class='btn' @click='uiState.modal.displayContent(`howToPlay`)'>HOW TO PLAY?</button>
        </div>
        <button class='modal-close' @click='close()'>&#10005;</button>
    </div>
</template>
<style scoped>
    @import '@/assets/styles/modal.css';
</style>
