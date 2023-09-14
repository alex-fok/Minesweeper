import { reactive } from 'vue'
import { GAMESTATUS } from '@/config'

export default reactive({
    id: '',
    roomId: -1,
    status: GAMESTATUS.UNDETERMINED,
    isPlayer: false,
    inviteCode: '',
    reset: function() {
      this.roomId = -1
      this.status = GAMESTATUS.UNDETERMINED
      this.isPlayer = false
      this.inviteCode = ''
    }
})
