import { reactive } from 'vue'
import { BOARDSETTING, GAMESTATUS } from '@/config'

type BlockView = {
    x: number,
    y: number,
    show: string,
    owner: string
}

type BlockInfo = {
    x: number,
    y: number,
    bType: number,
    value: number
}

type Player = {
    id: string,
    alias: string,
    score: number,
    isTurn: boolean,
    isOnline: boolean
}

const [BLANK, BOMB, NUMBER] = [0, 1, 2]

const board = new Array(BOARDSETTING.SIZE * BOARDSETTING.SIZE)
    .fill({})
    .map((_, i) => ({
        x: i % BOARDSETTING.SIZE,
        y: Math.floor(i / BOARDSETTING.SIZE),
        show: '',
        owner: ''
    })) as BlockView[]

const players: Record<string, Player> = {}
const bombColor: Record<string, string> = {}
export default reactive({
    id: "",
    roomId: -1,
    board,
    status: GAMESTATUS.NEW,
    resetBoard: function() {
        this.board = this.board.map((_, i) => ({
            ...this.board[i], ...{ show: '', owner: '' }
        }))
    },
    getDisplayVal: function (block: BlockInfo) : string {
        if (block['bType'] === NUMBER) return block['value'].toString()
        return block['bType'] === BOMB ? 'BO' : 'BL'
    },
    players,
    bombColor,
    bombsLeft: Number.MAX_SAFE_INTEGER,
    isGameOver: false,
    winner: ''
})
