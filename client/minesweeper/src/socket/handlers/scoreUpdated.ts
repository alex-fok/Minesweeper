import updateCounter from './shared/updateCounter'

export default (data: { player: number, opponent: number, bombsLeft: number }) => {
    const { player, opponent, bombsLeft } = data
    updateCounter(player, opponent, bombsLeft)
}
