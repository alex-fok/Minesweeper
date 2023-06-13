import { gameState } from '@/store'

export default (data: { client: string }) => {
    const { client } = data
    gameState.players[client].isOnline = false
}
