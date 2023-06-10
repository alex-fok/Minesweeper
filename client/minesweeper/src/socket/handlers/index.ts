import userIdHandler from './userId'
import roomIdHandler from './roomId'
import roomCreatedHandler from './roomCreated'
import gameStatHandler from './gameStat'
import scoreUpdatedHandler from './scoreUpdated'
import gameWonHandler from './gameWon'
import gameLostHandler from './gameLost'
import turnPassedHandler from './turnPassed'

import revealHandler from './reveal'

type handlerMap = {
    name: string,
    fn: (event: any) => void
}

const getAll = () : handlerMap[] => [
    { name: 'userId', fn: userIdHandler },
    { name: 'roomId', fn: roomIdHandler },
    { name: 'roomCreated', fn: roomCreatedHandler },
    { name: 'gameStat', fn: gameStatHandler },
    { name: 'scoreUpdate', fn: scoreUpdatedHandler },
    { name: 'gameWon', fn: gameWonHandler },
    { name: 'gameLost', fn: gameLostHandler },
    { name: 'turnPassed', fn: turnPassedHandler },
    { name: 'reveal', fn: revealHandler },
]

export default { getAll }
