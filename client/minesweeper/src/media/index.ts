import startUrl from '@/assets/media/start.mp3'
import winUrl from '@/assets/media/win.mp3'
import loseUrl from '@/assets/media/lose.mp3'
import scoreUrl from '@/assets/media/score.mp3'
import turnUrl from '@/assets/media/turn.mp3'

type SoundEffect = 'start' | 'win' | 'lose' | 'score' | 'turn'

let playing : SoundEffect | '' = ''

const audioMap:Record<string, HTMLAudioElement> = {
    start : new Audio(startUrl),
    win: new Audio(winUrl),
    lose: new Audio(loseUrl),
    score: new Audio(scoreUrl),
    turn: new Audio(turnUrl)
}

Object.getOwnPropertyNames(audioMap).forEach(se => {
    audioMap[se].addEventListener('ended', _ => {
        playing = ''
    })
})

export default {
    play: (se: SoundEffect) => {
        if (playing !== '') audioMap[playing].load()
        playing = se
        audioMap[se].play()
    }
}
