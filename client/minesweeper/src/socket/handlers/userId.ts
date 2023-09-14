import { roomState } from '@/store'

export default (data: { id: string }) => {
    if (roomState.id !== '') return
    const url = new URL(window.location.href)
    
    roomState.id = data.id
    url.searchParams.set('id', data.id)
    history.replaceState({}, '', url)
}
