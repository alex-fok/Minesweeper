<script setup lang='ts'>
defineProps({
    setContent: {
        type: Function,
        default: (content: string) => {}
    },
    close: {
        type: Function,
        default: () => {}
    }
})
</script>
<template>
    <div class='moda-row'>
        <div class='modal-item'>
            <span class='btn' @click='setContent(`create`)'>CREATE ROOM</span>
        </div>
        <div class='modal-item'>
            OR
        </div>
        <div class='modal-item'>
            <span class='btn' @click='setContent(`join`)'>JOIN ROOM</span>
        </div>
        <div class='modal-close' @click='close()'>&#10005;</div>
    </div>
</template>
<style scoped>
    @import '@/assets/modal.css';
</style>
