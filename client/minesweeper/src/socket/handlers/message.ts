import { uiState } from '@/store'

type Message = {
    message: string
}

export default (data: Message) => {
    const { message } = data
    uiState.modal.displayContent('message', message)
}
