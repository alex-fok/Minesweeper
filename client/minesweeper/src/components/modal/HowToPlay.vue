<script setup lang='ts'>
import { GAMESTATUS } from '@/config'
import { gameState, uiState } from '@/store'
import { computed, ref } from 'vue'

defineProps({
    close: {
        type: Function,
        default: () => {}
    }
})

const page = ref(1)
const maxPage = 4

const prevBtn = computed(() => page.value === 1 && gameState.status !== GAMESTATUS.NEW ? 'btn hidden' : 'btn')
const nextBtn = computed(() => page.value === maxPage ? 'btn hidden' : 'btn')

const nextPage = () => {
    if (page.value < maxPage)
        page.value += 1
}

const previousPage = () => {
    page.value <= 1 ?
        uiState.modal.displayContent('createOrJoin') :
        page.value -= 1
}
</script>
<template>
<div class='modal-row reverse'><div class='modal-close' @click='close()'>&#10005;</div></div>
<div class='modal-row'>
    <div class='modal-item'>
        <div class='instruction' v-if='page === 1'>
            <p>
                <span class='subtitle'>Objective:</span>
                The goal of the game is to reveal more mines than your opponent. The player who finds more than half of them first wins the game.
            </p>
        </div>
        <div class='instruction' v-else-if='page === 2'>
            <p>
                <span class='subtitle'>Mine blocks:</span>
                Player who reveals a mine block scores a point. Player's turn remains until revealing of a non-mine block.
            </p>
            <img src='@/assets/media/mine_block.png' />
        </div>
        <div class='instruction' v-else-if='page === 3'>
            <p>
                <span class='subtitle'>Number blocks:</span>
                Each number represents the number of mines in its neighboring blocks (horizontal, vertical and diagonal).
            </p>
            
            <img src='@/assets/media/num_block.png'/>
            <p>
                For example, the center block above has the number 1, which means there exists exactly 1 mine out of all 8 neighboring blocks.
            </p>
        </div>
        <div class='instruction' v-else>
            <p>
                <span class='subtitle'>Blank blocks:</span>
                If a block is blank, that simply means there exists no mine in its neighboring blocks.
            </p>
            <p>
                Revealing a blank block also reveals all other connected blank blocks (if any), as well as the number blocks surrounding them.
            </p>
            <img src='@/assets/media/blank_block.png' />
        </div>
    </div>
    
</div>
<div class='modal-row'>
    <div class='modal-item'>
        <button :class='prevBtn' @click='previousPage()'>
            {{gameState.status === GAMESTATUS.NEW && page === 1 ? 'BACK' : '&#8592; PREV' }}
        </button>
    </div>
    <div class='modal-item'>
        <button :class='nextBtn' @click='nextPage()'>NEXT &#8594;</button>
    </div>
</div>
</template>
<style scoped>
.instrution p {
    margin-top: 0;
    user-select: text;
    word-break: keep-all;
}
.instruction img {
    display: block;
    margin: 0 auto 1rem;
}
.subtitle {
    font-weight: 700;
}
</style>
