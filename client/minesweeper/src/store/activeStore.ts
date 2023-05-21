import { reactive } from 'vue'

export default reactive({
    active: true,
    updateActivity: function(bool: boolean) {
        this.active = bool
    }
})
