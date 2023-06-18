<script setup lang='ts'>
import { gameState } from '@/store'
import Flag from './icon/flag.vue'
import { uiState } from '@/store'
import { computed } from 'vue'
const props = defineProps({
    reveal: {
        type: Function,
        default: () => {}
    },
    show: {
        type: String,
        default: ''
    },
    owner: {
        type: String,
        default: '' 
    }
})

const { isPlayer, players, id } = gameState

const emitReveal = () => {
    props.reveal()
}

const isSelectable = computed(() => isPlayer && players[id].isTurn)
const getNumClass = (num: number) => `num-${num}`
</script>
<template>
    <div
    :class='`block${show !== `` ? ` revealed` : ``}${isSelectable ? ` selectable` : ``}`'
    @click='emitReveal'
    >
        <Flag
            v-if='show === `BO`'
            size='3vh'
            :fill='uiState.bombColor[owner]'
        />
        <span
            v-else-if='!Number.isNaN(parseInt(show))'
            :class= 'getNumClass(parseInt(show))'
        >{{ parseInt(show) }}</span>
    </div>
</template>
<style scoped>
    .block {
        font-family:'Franklin Gothic Medium', 'Arial Narrow', Arial, sans-serif;
        background-color:#444444;
        width: 3vh;
        height: 3vh;
        line-height: 3vh;
        vertical-align: middle;
        text-align: center;
        user-select: none;
    }
    .selectable:hover {
        outline: .1rem solid #CCCCCC;
        cursor: pointer;
    }
    .revealed {
        background-color: #343434;
    }
    .num-1 {
        color: cornflowerblue;
    }
    .num-2 {
        color: darkcyan;
    }
    .num-3 {
        color: crimson;
    }
    .num-4 {
        color: blueviolet;
    }
    .num-5 {
        color: gold;
    }
    .num-6 {
        color: turquoise;
    }
    .num-7 {
        color: pink;
    }
    .num-8 {
        color: gray;
    }
</style>
