<script setup lang='ts'>
defineProps({
    turnCount: Number 
})
</script>
<template>
    <div class='side-container'>
        <div class='player-list'>
        <div class='player-item'>
            <div class='player-name'><span class='curr hide'>>></span><p>Player 1 &#128081;</p></div>
            <div class='score'>31</div>
        </div>
        <div class='player-item selected'>
            <div class='player-name'><span class='curr'>>></span><p>MMMMMMMMMMMMMMMMMM</p></div>
            <div class='score'>28</div>
        </div>
        </div>
    </div>
</template>
<style scoped>
    .side-container {
        margin-top: 0%;
        padding: 2rem 2rem;
        box-sizing: border-box;
        background-color: #222;
        width: 18rem;
        height: 100%;
        display: flex;
        flex-direction: column;
        color: #9F9F9F;
    }
    .player-list {
        display: flex;
        flex-direction: column;
        row-gap: 1rem;
        height: 65%;
    }
    .player-item {
        border: 1px solid #9F9F9F;
        border-radius: .5rem;
        flex-grow: 1;
        padding: 1rem 1rem;
        display: flex;
        flex-direction: column;
    }
    .player-item.selected {
        border-color: white;
        color: white;
    }
    .player-name {
        display: flex;
        flex-direction: row;
        font-size: 1.2rem;
        flex-grow: 1;
    }
    .player-name .curr {
        width: 2rem;
    }
    .player-name p {
        margin: 0;
        word-break: break-all;
    }
    .score {
        font-size: 6rem;
        text-align: center;
        margin: auto;
    }
    .hide {
        visibility: hidden;
    }
</style>
