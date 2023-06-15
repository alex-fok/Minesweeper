import { gameState } from '@/store'

type PlayerAlias = {
    client: string,
    alias: string
}

export default (data: PlayerAlias) => {
    const { client, alias } = data
    if (gameState.players[client]) {
        console.log(client, ' assigned alias:', alias)
        gameState.players[client].alias = alias
    }
}
