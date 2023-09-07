import { reactive } from 'vue'

type Content =
    'create' | 'join' | 'roomSearch' |
    'createOrJoin' | 'gameEnded' | 'invited' |
    'message' | 'howToPlay' | 'passcode' |
    'playerAlias' | 'waitingRoom'

const playerColor: Record<string, string> = {}

export default reactive({
    modal: {
        isActive: false,
        content: 'createOrJoin',
        message: '',
        displayContent: function(v: Content, msg: string ='') {
            this.content = v,
            this.isActive = true
            this.message = msg
        }
    },
    playerColor: playerColor
})
