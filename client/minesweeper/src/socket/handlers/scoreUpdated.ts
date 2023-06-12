import updateCounter from './shared/updateCounter'
import { gameState } from '@/store'

export default (data: { score: Record<string, number>, bombsLeft: number }) => {
    const { score, bombsLeft } = data
    updateCounter(score, bombsLeft)
}
