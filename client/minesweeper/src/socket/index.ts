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

    const reconnect = () => {
        console.log('Trying to reconnect user')
        socket.send(JSON.stringify({
            name: 'reconnect',
            content: JSON.stringify({ userId, roomId })
        }))
    }

    const confirmInviteCode = () => {
        gameState.status = INVITED
        socket.send(JSON.stringify({
            name: 'inviteCode',
            content: JSON.stringify({id: invitation})
        }))
    }

    gameState.id = userId || ''
    if (userId) {
        reconnect()
    } else if (invitation) {
       confirmInviteCode()
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
