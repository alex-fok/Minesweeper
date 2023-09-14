import { GAMESTATUS } from '@/config'
import { roomState, reset } from '@/store'

type RoomId = {
    id: number,
    inviteCode: string
}

const { WAITING_JOIN } = GAMESTATUS

export default (data: RoomId) => {
    reset()
    const { id, inviteCode } = data
    roomState.roomId = id
    roomState.inviteCode = inviteCode
    if (![undefined, -1].includes(id)) {
        document.title = `#${id} Minesweeper`
    }
    // Change search query
    const url = new URL(window.location.href)
    url.searchParams.set('room', data.id.toString())
    history.replaceState({}, '', url)

    roomState.status = WAITING_JOIN
}
