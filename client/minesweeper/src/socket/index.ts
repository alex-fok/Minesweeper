import { GAMESTATUS } from '@/config'
import { gameState, uiState } from '@/store'
import handlers from './handlers'

const { NEW, INVITED } = GAMESTATUS

// For prod: window.location.hostname
// For dev : ws://localhost:8080/ws
const socket = new WebSocket(import.meta.env.VITE_SERVER ? import.meta.env.VITE_SERVER : 'ws://' + window.location.hostname + '/ws')
const socketEvents: Record<string, (event: any)=>void> = {}

handlers.getAll().forEach(handler => {
    socketEvents[handler.name] = handler.fn
})

socket.addEventListener('open', _ => {
    const url = new URL(window.location.href)
    const roomId = url.searchParams.get('room')
    const userId = url.searchParams.get('id')
    const invitation = url.searchParams.get('join')

    if (userId) {
        gameState.id = userId
        if (roomId) gameState.roomId = parseInt(roomId)

        socket.send(JSON.stringify({
            name: 'reconnect',
            content: JSON.stringify({ userId, roomId })
        }))
    } else if (invitation) {
        gameState.status = INVITED
        socket.send(JSON.stringify({
            name: 'inviteCode',
            content: JSON.stringify({id: invitation})
        }))
    } else {
        uiState.modal.displayContent('createOrJoin')
        gameState.status = NEW
    }
})

socket.addEventListener('message', event => {
    const { name, content } = JSON.parse(event.data)
    console.log('name: ', name, '\ncontent: ', content)
   if (!socketEvents[name]) return
   socketEvents[name](JSON.parse(content))
})

export default socket
