import { gameState } from '@/store'

type mousePosition = {
    id: string,
    position: number
}

export default (data: mousePosition) => {
    const { id, position } = data
    if (gameState.players[id])
        gameState.players[id].cursor = position 
}
