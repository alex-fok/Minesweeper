
const socket = new WebSocket("ws://localhost:8080/ws")

const socketEvents: Record<string, (event: any)=>void> = {}

socket.addEventListener("message", (event) => {
   const eventName =  "reveal"
   socketEvents[eventName](event)
})

const addSocketEvent = (name: string, fn: (event: any)=>void) => {
    socketEvents[name] = fn
}
export { socket, addSocketEvent }
