<template>
  <div id="app">
    <img class="bg-img" src="https://img-pool.zxysy.net/nsbg.jpg">
    <Barrage></Barrage>
    <div class="qr-text"><p>扫码发弹幕</p></div>
    <div class="qr-code">
      <img class="qr-img" src="https://img-pool.zxysy.net/qrCode.png">
    </div>
  </div>
</template>

<script>

import Barrage from "./components/Barrage"

export default {
  name: 'app',
  components: {
    Barrage: Barrage
  }
}
</script>

<style>
#app {
  font-family: 'Avenir', Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  /* text-align: center; */
  /* color: #2c3e50; */
  margin-top: 60px;
}

.bg-img {
  width: 100%;
  height: 100%;
  position: fixed;
  top: 0;
  left: 0;
  z-index: -2;
}

.qr-code {
  overflow: hidden;
  position: fixed;
  width: 300px;
  height: 300px;
  right: 30px;
  bottom: 30px;
  background: rgba(255, 255, 255, 95%);
  border-radius: 15px;
  z-index: -1;
}

.qr-img {
  width: 100%;
  height: 100%;
}

.qr-text {

  font-family: 'Avenir', Helvetica, Arial, sans-serif;
  color: #2c3e50;

  position: fixed;
  width: 300px;
  height: 40px;
  right: 30px;
  bottom: 350px;
  background: rgba(255, 255, 255, 95%);
  border-radius: 15px;
  z-index: -1;

  display: flex;
  justify-content: center;
  align-items: center;
}
</style>
