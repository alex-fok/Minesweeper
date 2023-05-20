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
        show: ""
    })) as Block[]

const resetBoard = () => {
    store.board = board.map((_, i) => ({
        ...board[i], ...{ show: '' }
    }))
}

export const store = reactive({
    board,
    resetBoard
})
