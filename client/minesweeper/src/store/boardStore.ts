import { reactive } from 'vue'
import { BOARDSETTING } from '@/config'

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
    resetBoard: function() {
        this.board = this.board.map((_, i) => ({
        ...this.board[i], ...{ show: '' }
    }))
    }
})
