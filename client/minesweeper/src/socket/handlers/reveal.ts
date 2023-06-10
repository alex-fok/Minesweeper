import { gameState } from '@/store'
import { BOARDSETTING } from '@/config'

type BlockInfo = {
    x: number,
    y: number,
    bType: number,
    value: number
}

export default (data: {blocks:BlockInfo[]}) => {
    const { blocks } = data
    
    blocks.forEach(block => {
        gameState.board[BOARDSETTING.SIZE * block.y + block.x].show = gameState.getDisplayVal(block)
    })
}
