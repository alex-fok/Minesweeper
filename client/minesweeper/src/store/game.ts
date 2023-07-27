import { reactive } from 'vue'
import { GAMESTATUS, BLOCKTYPE, BOARDSETTING } from '@/config'

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
    isOnline: boolean,
    cursor: number
}

const { UNDETERMINED } = GAMESTATUS
const { BOMB, NUMBER } = BLOCKTYPE

const board: BlockView[] = [];

const players: Record<string, Player> = {}

export default reactive({
    id: '',
    roomId: -1,
    board,
    status: UNDETERMINED,
    isPlayer: false,
    boardConfig: {
        bomb: BOARDSETTING.BOMB.NORMAL,
        size: BOARDSETTING.SIZE.MEDIUM
    },
    initBoard: function() {
        this.board = new Array(this.boardConfig.size * this.boardConfig.size)
        .fill({})
        .map((_, i) => ({
            x: i % this.boardConfig.size,
            y: Math.floor(i / this.boardConfig.size),
            show: '',
            owner: ''
        })) as BlockView[]
    },
    resetBoard: function() {
        this.board = this.board.map((_, i) => ({
            ...this.board[i], ...{ show: '', owner: '' }
        }))
    },
    getDisplayVal: function (block: BlockInfo) : string {
        if (block['bType'] === NUMBER) return block['value'].toString()
        return block['bType'] === BOMB ? 'BO' : 'BL'
    },
    players: players,
    bombsLeft: Number.MAX_SAFE_INTEGER,
    inviteCode: '',
    winner: ''
})
