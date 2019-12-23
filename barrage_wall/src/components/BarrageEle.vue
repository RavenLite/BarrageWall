<template>
    <div class="barrage-element-root" :style="`animation: barrage ${lengthFunc(text.length)}s linear 0s 1;`">
        <p class="barrage-text" :style="`color: ${color}`">  {{text}}</p>
    </div>
</template>

<script>
export default {
    props: {
        text: String,
        color: String
    },
    methods: {
        lengthFunc: function(len) {
            return 40/Math.sqrt(len)
        }
    }
}
</script>

<style>
.barrage-text {
    word-break:keep-all;
    white-space:nowrap;
    font-size: 60px;
    -webkit-text-stroke-width: 1.5px;
    -webkit-text-stroke-color: black;
    font-weight: 800;
}

.barrage-element-root {
    height: 65px;
    position: absolute;
    left: -3000px;
}

@keyframes barrage {
    from {
        left: 130%;
        transform: translateX(0);
    } to {
        left: 0;
        transform: translateX(-100%);
    }
}
</style>