import { gameState } from '@/store'

type PlayerOnline = {
    player: string
    isOnline: boolean
}

export default (data: PlayerOnline) => {
    const { player, isOnline } = data
    gameState.players[player].isOnline = isOnline 
}
