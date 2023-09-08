import { publicState } from '@/store'

type roomIds = {
    ids: number[]
}

export default (data: roomIds) => {
    console.log(data)
    publicState.roomIds = data.ids
}
