<template>
    <div class='app-container'>
        <slot name='header'></slot>
        <div class='main-container'>
            <slot></slot>
        </div>
    </div>
</template>
<style scoped>
    .app-container {
        display: flex;
        flex-direction: column;
        height: 100vh;
    }
    .main-container {
        flex-grow: 1;
        display: flex;
        flex-direction: row;
        height: 100%;
    }
</style>
