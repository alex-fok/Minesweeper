import { gameState } from '@/store'
import { GAMESTATUS } from '@/config'

const { WAITING_JOIN } = GAMESTATUS

export default () => {
    gameState.status = WAITING_JOIN
    gameState.resetBoard()
}
