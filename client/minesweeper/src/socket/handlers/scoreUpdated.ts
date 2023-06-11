import updateCounter from './shared/updateCounter'
import { gameState } from '@/store'

export default (data: { score: Record<string, number>, bombsLeft: number }) => {
    const { score, bombsLeft } = data
    const [player, opponent] = [score[gameState.player.id], score[gameState.opponent.id]]
    updateCounter(player, opponent, bombsLeft)
}
