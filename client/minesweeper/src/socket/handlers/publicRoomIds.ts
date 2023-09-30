import { publicState } from '@/store'

type publicRoomInfo = {
    id: number,
    capacity: number,
    users: number
}

type rooms = {
    rooms: publicRoomInfo[]
}

export default (data: rooms) => {
    publicState.rooms = data.rooms
}
