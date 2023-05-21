import { reactive } from 'vue'
import { boardSetting } from '@/config'

type Block = {
    x: number,
    y: number,
    show: string
}

const board = new Array(boardSetting.SIZE * boardSetting.SIZE)
    .fill({})
    .map((_, i) => ({
        x: i % boardSetting.SIZE,
        y: Math.floor(i / boardSetting.SIZE),
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
