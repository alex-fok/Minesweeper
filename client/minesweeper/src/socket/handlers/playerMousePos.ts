import { gameState } from '@/store'

type mousePosition = {
    id: string,
    position: number
}

export default (data: mousePosition) => {
    const { id, position } = data
    console.log(gameState.players)
    console.log('playerId: ', id)
    if (gameState.players[id])
        gameState.players[id].cursor = position 
}
