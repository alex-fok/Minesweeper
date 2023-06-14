import updateCounter from './shared/updateCounter'
import media from '@/media'

export default (data: { score: Record<string, number>, bombsLeft: number }) => {
    const { score, bombsLeft } = data
    updateCounter(score, bombsLeft)

    media.play('score')
}
