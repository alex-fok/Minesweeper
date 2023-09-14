import { reactive } from 'vue'
import { BLOCKTYPE, BOARDSETTING } from '@/config'

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
    isReady: boolean,
    cursor: number
}

const { BOMB, NUMBER } = BLOCKTYPE

const board: BlockView[] = [];

const players: Record<string, Player> = {}
export default reactive({
    board,
    boardConfig: {
        bomb: BOARDSETTING.BOMB.NORMAL,
        size: BOARDSETTING.SIZE.MEDIUM
    },
    players: players,
    capacity: 0,
    bombsLeft: Number.MAX_SAFE_INTEGER,
    timeLimit: 0,
    lastPlayed: {
        x: -1,
        y: -1,
        owner: '',
        timestamp: 0
    },
    winner: '',
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
    reset: function() {
        this.resetBoard()
        this.players = {}
        this.capacity = 0
        this.bombsLeft = Number.MAX_SAFE_INTEGER
        this.winner = ''
    },
})
