import gameEndedHandler from './gameEnded'
import gameStatHandler from './gameStat'
import messageHandler from './message'
import passcodeHandler from './passcode'
import playerAliasHandler from './playerAlias'
import playerMousePos from './playerMousePos'
import playerOnlineHandler from './playerOnline'
import playerReadyHandler from './playerReady'
import publicRoomIdsHandler from './publicRoomIds'
import reconnFailedHandler from './reconnFailed'
import revealHandler from './reveal'
import roomCreatedHandler from './roomCreated'
import roomInfoHandler from './roomInfo'
import scoreUpdatedHandler from './scoreUpdated'
import turnPassedHandler from './turnPassed'
import userIdHandler from './userId'
import waitingRoomHandler from './waitingRoom'

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
    { name: 'playerOnline', fn: playerOnlineHandler },
    { name: 'playerReady', fn: playerReadyHandler },
    { name: 'publicRoomIds', fn: publicRoomIdsHandler },
    { name: 'reconnFailed', fn: reconnFailedHandler },
    { name: 'reveal', fn: revealHandler },
    { name: 'roomCreated', fn: roomCreatedHandler },
    { name: 'roomInfo', fn: roomInfoHandler },
    { name: 'scoreUpdated', fn: scoreUpdatedHandler },
    { name: 'turnPassed', fn: turnPassedHandler },
    { name: 'userId', fn: userIdHandler },
    { name: 'waitingRoom', fn: waitingRoomHandler }
]

export default { getAll }
