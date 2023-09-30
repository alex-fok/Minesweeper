import { reactive } from 'vue';

type publicRooms = {
    id: number,
    capacity: number
    users: number
}

const rooms: publicRooms[] = []

export default reactive({
    rooms
})
