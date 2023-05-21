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
        justify-content: center;
        align-items: center;
        column-gap: 2rem;
        height: 100%;
    }
</style>
