<script setup lang='ts'>
import { computed } from 'vue'
import { BOARDSETTING } from '@/config'
import { gameState } from '@/store'
import Flag from './icon/flag.vue'
import { uiState } from '@/store'

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
    },
    playerHovering: {
        type: String || undefined,
        default: undefined
    },
    updateMousePosition: {
        type: Function,
        default: () => {}
    }
})

const { id } = gameState

const emitReveal = () => {
    props.reveal()
}
const emitPositionUpdate = () => {
    props.updateMousePosition()
}

const getNumClass = (num: number) => `num-${num}`
const isSelectable = computed(() => gameState.isPlayer && gameState.players[id].isTurn)
const isShrunk = computed(() => gameState.boardConfig.size === BOARDSETTING.SIZE.LARGE)
</script>
<template>
    <div
        :class='`block${show !== `` ? ` revealed` : ``}${isSelectable ? ` selectable` : ``}${isShrunk ? ` shrunk`: ``}`'
        :style='{outline: `${playerHovering ? `.1rem solid ` + uiState.playerColor[playerHovering] : ``}`}'
        @click='emitReveal'
        @mouseover='emitPositionUpdate'
    >
        <Flag
            v-if='show === `BO`'
            :size='isShrunk ? `2vh` : `3vh`'
            :fill='uiState.playerColor[owner]'
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
    .block.shrunk {
        width: 2vh;
        height: 2vh;
        line-height: 2vh;
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
