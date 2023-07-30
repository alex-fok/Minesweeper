import gameEndedHandler from './gameEnded'
import gameStatHandler from './gameStat'
import messageHandler from './message'
import passcodeHandler from './passcode'
import playerAliasHandler from './playerAlias'
import playerMousePos from './playerMousePos'
import playerOfflineHandler from './playerOffline'
import playerOnlineHandler from './playerOnline'
import reconnFailedHandler from './reconnFailed'
import revealHandler from './reveal'
import roomCreatedHandler from './roomCreated'
import roomInfoHandler from './roomInfo'
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
    { name: 'passcode', fn: passcodeHandler },
    { name: 'playerAlias', fn: playerAliasHandler },
    { name: 'playerMousePos', fn: playerMousePos },
    { name: 'playerOffline', fn: playerOfflineHandler },
    { name: 'playerOnline', fn: playerOnlineHandler },
    { name: 'reconnFailed', fn: reconnFailedHandler },
    { name: 'reveal', fn: revealHandler },
    { name: 'roomCreated', fn: roomCreatedHandler },
    { name: 'roomInfo', fn: roomInfoHandler },
    { name: 'scoreUpdated', fn: scoreUpdatedHandler },
    { name: 'turnPassed', fn: turnPassedHandler },
    { name: 'userId', fn: userIdHandler },
]

export default { getAll }
