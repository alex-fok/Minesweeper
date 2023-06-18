import gameEndedHandler from './gameEnded'
import gameStatHandler from './gameStat'
import messageHandler from './message'
import playerAliasHandler from './playerAlias'
import playerOfflineHandler from './playerOffline'
import playerOnlineHandler from './playerOnline'
import reconnFailedHandler from './reconnFailed'
import revealHandler from './reveal'
import roomCreatedHandler from './roomCreated'
import roomIdHandler from './roomId'
import scoreUpdatedHandler from './scoreUpdated'
import turnPassedHandler from './turnPassed'
import userIdHandler from './userId'

type handlerMap = {
    name: string,
    fn: (event: any) => void
}

const getAll = () : handlerMap[] => [
    { name: 'gameEnded', fn: gameEndedHandler },
    { name: 'gameStat', fn: gameStatHandler },
    { name: 'message', fn: messageHandler },
    { name: 'playerAlias', fn: playerAliasHandler },
    { name: 'playerOffline', fn: playerOfflineHandler },
    { name: 'playerOnline', fn: playerOnlineHandler },
    { name: 'reconnFailed', fn: reconnFailedHandler },
    { name: 'reveal', fn: revealHandler },
    { name: 'roomCreated', fn: roomCreatedHandler },
    { name: 'roomId', fn: roomIdHandler },
    { name: 'scoreUpdated', fn: scoreUpdatedHandler },
    { name: 'turnPassed', fn: turnPassedHandler },
    { name: 'userId', fn: userIdHandler },
]

export default { getAll }
