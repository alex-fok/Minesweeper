import { gameState } from '@/store'
import { GAMESTATUS } from '@/config'

const { PLAYING, WAITING_TURN } = GAMESTATUS

export default () => {
    gameState.status = gameState.status === PLAYING ? WAITING_TURN : PLAYING
}
