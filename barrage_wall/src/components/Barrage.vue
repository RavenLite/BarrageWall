<template>
    <div id="barrage-root">
        <div class="barrage-line-container" :style="`top:${topFunc(barrage.index)}px`"  v-for="barrage in barrageList" :key="barrage.index">
            <BarrageEle :color="barrage.color" :text="barrage.text"></BarrageEle>
        </div>
    </div>
</template>

<script>
import BarrageEle from './BarrageEle'

export default {
    components: {
        BarrageEle
    },
    data() {
        return {
            barrageCount: 1,
            debugText: "",
            lineCount: 10,
            barrageList: []
        }
    },
    watch: {
        barrageList: function(val) {
            if (val.length > 300) {
                this.barrageList = []
            }
        }
    },
    created() {
        let that= this;
        let ws = new WebSocket('wss://barrage.neuyan.com/ws?s=071a7756-6032-49e1-b9c8-35b1ca8df397');
        ws.onmessage = function(event) {
            window.console.log(JSON.parse(JSON.parse(event.data).Message).text)
            if (JSON.parse(JSON.parse(event.data).Message).text && JSON.parse(JSON.parse(event.data).Message).text.length>40) {
                return
            } 
            that.addBarrage(JSON.parse(JSON.parse(event.data).Message).text)
        }
        ws.onclose = function() {
            ws = new WebSocket('wss://barrage.neuyan.com/ws?s=071a7756-6032-49e1-b9c8-35b1ca8df397')
        }
        // setInterval(()=>this.addBarrage(Math.random().toString(36).substr(2)), 1000)
    },
    methods: {
        addBarrage: function(text) {
            let barrageObj = {text: text, color: this.color(this.barrageCount), index: this.barrageCount++}
            this.barrageList.push(barrageObj)
        },
        topFunc: function(index) {
            return (((index*index)+index*2)%8)*65
        },
        color: function(index) {
            const colors = ['#FFFF00','#FF8500','#00FFFF', '#FF0000','#00FF00','#0000FF','#E600FF', '#000']
            return colors[index%colors.length]
        } 
    }
}
</script>

<style>
html {
    overflow:hidden;
}

#barrage-root {
    height: 100%;
}

.barrage-line-container {
    overflow: inherit;
    position: relative;
    width: 80%;
}

</style>