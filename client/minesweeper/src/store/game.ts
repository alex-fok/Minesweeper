import { reactive } from 'vue'
import { BOARDSETTING, GAMESTATUS } from '@/config'

type Block = {
    x: number,
    y: number,
    show: string
}

const board = new Array(BOARDSETTING.SIZE * BOARDSETTING.SIZE)
    .fill({})
    .map((_, i) => ({
        x: i % BOARDSETTING.SIZE,
        y: Math.floor(i / BOARDSETTING.SIZE),
        show: ''
    })) as Block[]

export default reactive({
    board,
    status: GAMESTATUS.NEW,
    resetBoard: function() {
        this.board = this.board.map((_, i) => ({
            ...this.board[i], ...{ show: '' }
        }))
    },
    player: {
        name: 'Anonymous',
        score: 0,
    },
    opponent: {
        name: 'Anonymous',
        score: 0,
    },
    bombsLeft: Number.MAX_SAFE_INTEGER,
    isGameOver: false,
})
