import { gameState } from '@/store'
export default (score: Record<string, number>, bombsLeft: number) => {
    // gameState.players[gameState.].score = playerScore
    // gameState.opponent.score = opponentScore
    const playerIds = Object.keys(score)
    playerIds.forEach(id => {
        gameState.players[id].score = score[id]
    })
    gameState.bombsLeft = bombsLeft
}
