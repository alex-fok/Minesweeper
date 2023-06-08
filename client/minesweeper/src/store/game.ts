import { reactive } from 'vue'
import { BOARDSETTING, GAMESTATUS } from '@/config'

type BlockView = {
    x: number,
    y: number,
    show: string
}

type BlockInfo = {
    x: number,
    y: number,
    bType: number,
    value: number
}

const [BLANK, BOMB, NUMBER] = [0, 1, 2]

const board = new Array(BOARDSETTING.SIZE * BOARDSETTING.SIZE)
    .fill({})
    .map((_, i) => ({
        x: i % BOARDSETTING.SIZE,
        y: Math.floor(i / BOARDSETTING.SIZE),
        show: ''
    })) as BlockView[]

export default reactive({
    board,
    status: GAMESTATUS.NEW,
    resetBoard: function() {
        this.board = this.board.map((_, i) => ({
            ...this.board[i], ...{ show: '' }
        }))
    },
    getDisplayVal: function (block: BlockInfo) : string {
        if (block['bType'] === NUMBER) return block['value'].toString()
        return block['bType'] === BOMB ? 'BO' : 'BL'
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
