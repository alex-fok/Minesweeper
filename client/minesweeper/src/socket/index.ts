import { GAMESTATUS } from '@/config'
import { roomState, uiState } from '@/store'
import handlers from './handlers'

const { NEW, INVITED } = GAMESTATUS

class CustomWebSocket extends WebSocket {
    constructor(url: string | URL, protocols?: string | string[] | undefined) {
        super(url, protocols)
    }
    emit(name: string, content: Object) {
        this.send(JSON.stringify({
            name,
            content: JSON.stringify(content)
        }))
    }
}

// For prod: window.location.hostname
// For dev : ws://localhost:8080/ws
const socket = new CustomWebSocket(import.meta.env.VITE_SERVER ? import.meta.env.VITE_SERVER : 'ws://' + window.location.hostname + '/ws')
const socketEvents: Record<string, (event: any)=>void> = {}

handlers.getAll().forEach(handler => {
    socketEvents[handler.name] = handler.fn
})

socket.addEventListener('open', () => {
    const url = new URL(window.location.href)
    const roomId = url.searchParams.get('room')
    const userId = url.searchParams.get('id')
    const invitation = url.searchParams.get('join')

    const reconnect = () => {
        console.log('Trying to reconnect user')
        socket.emit('reconnect', { userId, roomId })
    }

    const confirmInviteCode = () => {
        roomState.status = INVITED
        socket.emit('inviteCode', { id: invitation })
    }

    roomState.id = userId || ''
    if (userId) {
        reconnect()
    } else if (invitation) {
       confirmInviteCode()
    } else {
        uiState.modal.displayContent('createOrJoin')
        roomState.status = NEW
    }
})

socket.addEventListener('message', event => {
    const { name, content } = JSON.parse(event.data)
    console.log('name: ', name, '\ncontent: ', content)
   if (!socketEvents[name]) return
   socketEvents[name](JSON.parse(content))
})

export default socket
