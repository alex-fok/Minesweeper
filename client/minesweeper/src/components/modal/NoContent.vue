<script setup lang='ts'>
defineProps({
    close: {
        type: Function,
        default: () => {}
    }
})
</script>
<template>
    <div class='modal-row'>
        No Content
        <div class='modal-close' @click='close()'>&#10005;</div>
    </div>
</template>
