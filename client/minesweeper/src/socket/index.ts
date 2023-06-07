// For prod: window.local.hostname
// For dev : ws://localhost:8080/ws
const socket = new WebSocket(import.meta.env.VITE_SERVER ? import.meta.env.VITE_SERVER : 'ws://' + window.location.hostname + '/ws')

const socketEvents: Record<string, (event: any)=>void> = {}
socket.addEventListener('open', _ => {
    socket.send(JSON.stringify({name: 'newRoom'})) 
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
    url.searchParams.set('id', id)
    history.replaceState({}, "", url)
})

export { socket, addSocketEventHandler }
