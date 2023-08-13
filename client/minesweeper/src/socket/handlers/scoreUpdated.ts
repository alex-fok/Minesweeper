import { gameState } from '@/store'
import updateCounter from './shared/updateCounter'
import media from '@/media'

type Counter = {
    score: Record<string, number>,
    bombsLeft: number,
    lastPlayed: string
}

export default (data: Counter) => {
    const { score, bombsLeft, lastPlayed } = data
    updateCounter(score, bombsLeft)
    gameState.lastPlayed.timestamp = Date.parse(lastPlayed)
    media.play('score')
}
