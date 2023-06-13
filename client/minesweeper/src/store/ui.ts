import { reactive } from 'vue'

type Content = 'create' | 'join' | 'createOrJoin' | 'gameEnded'

const bombColor: Record<string, string> = {}

export default reactive({
    modal: {
        isActive: true,
        content: 'createOrJoin',
        displayContent: function(v: Content) {
            this.content = v,
            this.isActive = true
        }
    },
    bombColor
})
