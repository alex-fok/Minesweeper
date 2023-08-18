import { gameState } from '@/store'
export default (score: Record<string, number>, bombsLeft: number) => {
    const playerIds = Object.keys(score)
    playerIds.forEach(id => {
        gameState.players[id].score = score[id]
    })
    gameState.bombsLeft = bombsLeft
}
