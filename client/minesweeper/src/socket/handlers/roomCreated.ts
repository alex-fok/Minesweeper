import { gameState } from '@/store'
import { GAMESTATUS } from '@/config'

type RoomCreated = {
    inviteCode: string
}

const { WAITING_JOIN } = GAMESTATUS

export default (data: RoomCreated) => {
    const { inviteCode } = data
    gameState.status = WAITING_JOIN
    gameState.inviteCode = inviteCode
    gameState.resetBoard()
}
