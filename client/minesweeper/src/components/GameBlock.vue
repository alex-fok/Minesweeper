<script setup lang='ts'>
import { computed } from 'vue'
import { BOARDSETTING } from '@/config'
import { gameState, roomState, uiState } from '@/store'
import Flag from './icon/FlagIcon.vue'

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
    isLastPlayed: {
        type: Boolean,
        default: false
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

const emitReveal = () => {
    props.reveal()
}
const emitPositionUpdate = () => {
    props.updateMousePosition()
}

const getNumClass = (num: number) => `num-${num}`
const isSelectable = computed(() => roomState.isPlayer && gameState.players[roomState.id].isTurn)
const isShrunk = computed(() => gameState.boardConfig.size === BOARDSETTING.SIZE.LARGE)
const outline = computed(() => {
    // Last Hand
    if (props.isLastPlayed) {
        const playerColor = uiState.playerColor[gameState.lastPlayed.owner]
        return `2px dashed rgba(${playerColor})`
    }
    // Hovering
    if (!props.playerHovering) return ''
    const { isTurn } = gameState.players[props.playerHovering]
    const playerColor = uiState.playerColor[props.playerHovering]

    return `2px solid ${isTurn ? `rgba(${playerColor})` : `rgba(${playerColor}, .5)`}`
})
</script>
<template>
    <div
        :class='`block${show !== `` ? ` revealed` : ``}${isSelectable ? ` selectable` : ``}${isShrunk ? ` shrunk`: ``}`'
        :style='{outline: outline}'
        @click='emitReveal'
        @mouseover='emitPositionUpdate'
    >
        <Flag
            v-if='show === `BO`'
            :size='isShrunk ? `2vh` : `3vh`'
            :color='`rgba(${uiState.playerColor[owner]})`'
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
        outline: 2px solid #FFFFFF !important;
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
