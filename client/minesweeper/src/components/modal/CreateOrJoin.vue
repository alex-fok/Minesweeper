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
            <span class='btn' @click='uiState.modal.displayContent(`create`)'>CREATE ROOM</span>
        </div>
        <div class='modal-item'>
            OR
        </div>
        <div class='modal-item'>
            <span class='btn' @click='uiState.modal.displayContent(`join`)'>JOIN ROOM</span>
        </div>
        <div class='modal-close' @click='close()'>&#10005;</div>
    </div>
</template>
<style scoped>
    @import '@/assets/styles/modal.css';
</style>
