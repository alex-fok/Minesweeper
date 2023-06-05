<script setup lang='ts'>
const props = defineProps({
    reveal: {
        type: Function,
        default: () => {}
    },
    show: {
        type: String,
        default: ''
    }
})

const emitReveal = () => {
    props.reveal()
}
</script>
<template>
    <div
    :class='show === `` ? `block` : `block revealed`'
    @click='emitReveal'
    >
        {{ show }}
    </div>
</template>
<style scoped>
    .block {
        font-family:'Courier New', Courier, monospace;
        background-color:#444444;
        width: 3vh;
        height: 3vh;
        line-height: 3vh;
        vertical-align: middle;
        text-align: center;
        user-select: none;
    }
    .revealed {
        background-color: #343434;
    }
</style>
