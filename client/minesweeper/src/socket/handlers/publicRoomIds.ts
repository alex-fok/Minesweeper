import { publicState } from '@/store'

type roomIds = {
    ids: number[]
}

export default (data: roomIds) => {
    publicState.roomIds = data.ids
}
