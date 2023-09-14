import { gameState } from '@/store'

type PlayerAlias = {
    client: string,
    alias: string
}

export default (data: PlayerAlias) => {
    const { client, alias } = data
    if (gameState.players[client]) {
        gameState.players[client].alias = alias
    }
}
