
const socket = new WebSocket("ws://localhost:8080/ws")

const socketEvents: Record<string, (event: any)=>void> = {}
socket.addEventListener("open", event => {
    socket.send(JSON.stringify({name: 'newGame'})) 
})
socket.addEventListener("message", event => {
   const eventName =  "reveal"
   if (!socketEvents[eventName]) return
   socketEvents[eventName](event)
})

const addSocketEventHandler = (name: string, fn: (event: any)=>void) => {
    socketEvents[name] = fn
}
export { socket, addSocketEventHandler }
