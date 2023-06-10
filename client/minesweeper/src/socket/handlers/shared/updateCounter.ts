import { gameState } from '@/store'
export default (playerScore: number, opponentScore: number, bombsLeft: number) => {
    gameState.player.score = playerScore
    gameState.opponent.score = opponentScore
    gameState.bombsLeft = bombsLeft
}
