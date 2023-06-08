// For prod: window.location.hostname
// For dev : ws://localhost:8080/ws
const socket = new WebSocket(import.meta.env.VITE_SERVER ? import.meta.env.VITE_SERVER : 'ws://' + window.location.hostname + '/ws')

const socketEvents: Record<string, (event: any)=>void> = {}
socket.addEventListener('open', _ => {
    const url = new URL(window.location.href)
    const roomId = url.searchParams.get('room')
    const userId = url.searchParams.get('id')
    // console.log('open -> roomId:', roomId)
    // console.log('open -> userId:', userId)

    if (!userId) return
    socket.send(JSON.stringify({
        name: 'reconnect',
        content: JSON.stringify({ userId, roomId })
    })) 
})

socket.addEventListener('message', event => {
    const { name, content } = JSON.parse(event.data)
    console.log('name: ', name)
    console.log('content: ', content)
   if (!socketEvents[name]) return
   socketEvents[name](JSON.parse(content))
})

const addSocketEventHandler = (name: string, fn: (event: any)=>void) => {
    socketEvents[name] = fn
}

addSocketEventHandler('userId', (content: { id: string }) => {
    const { id } = content
    const url = new URL(window.location.href)
    
    if (url.searchParams.get('id')) return

    url.searchParams.set('id', id)
    history.replaceState({}, "", url)
})

export { socket, addSocketEventHandler }
