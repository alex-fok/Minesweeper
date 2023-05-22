<template>
    <slot name='header'></slot>
    <div class='main-container'>
        <slot></slot>
    </div>
</template>
<style scoped>
    .main-container {
        flex-grow: 1;
        display: flex;
        flex-direction: row;
        height: 100%;
    }
</style>
