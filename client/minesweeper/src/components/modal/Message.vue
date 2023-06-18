<script setup lang='ts'>
import { uiState } from '@/store';

defineProps({
    close: {
        type: Function,
        default: () => {},
    }
})
</script>
<template>
<div class='modal-row'>
    <div class='modal-item'>
        SERVER:
    </div>
    <div class='modal-item'>
        {{ uiState.modal.message }}
    </div>
    <div class='modal-item'>
        <div class='modal-close' @click='close()'>&#10005;</div>
    </div>
</div>

</template>
<style scoped>
    @import '@/assets/styles/modal';
</style>
