
const socket = new WebSocket("ws://localhost:8080/ws")

const socketEvents: Record<string, (event: any)=>void> = {}
socket.addEventListener("open", _ => {
    socket.send(JSON.stringify({name: 'newGame'})) 
})
socket.addEventListener("message", event => {
    const { name, content } = JSON.parse(event.data)
    // console.log('name: ', name)
    // console.log('content: ', content)
   if (!socketEvents[name]) return
   socketEvents[name](content)
})

const addSocketEventHandler = (name: string, fn: (event: any)=>void) => {
    socketEvents[name] = fn
}
export { socket, addSocketEventHandler }
