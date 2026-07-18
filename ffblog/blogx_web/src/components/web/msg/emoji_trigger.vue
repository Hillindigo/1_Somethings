<script setup lang="ts">

import {IconFaceSmileFill} from "@arco-design/web-vue/es/icon";
import {ref} from "vue";

const emits = defineEmits(["select"])

function emojiClick(e: Event){
  const target = e.target as HTMLElement
  if (target.tagName === "SPAN" && !target.classList.contains("img")){
    emits("select", "emoji", target.innerText)
    visible.value = false
    return
  }
  if (target.tagName === "IMG"){
    emits("select", "img", (target as HTMLImageElement).src)
    visible.value = false
    return
  }
  if (target.tagName === "SPAN" && target.classList.contains("img")){
    emits("select", "img", (target.childNodes[0] as HTMLImageElement).src)
    visible.value = false
    return
  }
}

const visible =ref(false)
const imageList = ref<string[]>([])
function initImages(){
  const files :Record<string, any>= import.meta.glob("@/assets/img/emoji/*.png", {eager: true})
  for (const filesKey in files) {
    const path = files[filesKey].default
    imageList.value.push(new URL(path, import.meta.url).href)
  }
}
initImages()



</script>

<template>
  <a-trigger trigger="click" v-model:popup-visible="visible" position="top" :unmount-on-close="false">
    <span><IconFaceSmileFill></IconFaceSmileFill></span>
    <template #content>
      <div class="emoji_trigger_content">
        <a-tabs :default-active-key="1" position="bottom">
          <a-tab-pane :key="1" title="默认表情">
            <div class="emoji_list" @click="emojiClick">
              <span class="item">😃</span>
              <span class="item">😄</span>
              <span class="item">😁</span>
              <span class="item">😆</span>
              <span class="item">😅</span>
              <span class="item">😂</span>
              <span class="item">😉</span>
              <span class="item">😊</span>
              <span class="item">😇</span>
              <span class="item">😍</span>
              <span class="item">😘</span>
              <span class="item">😚</span>
              <span class="item">😋</span>
              <span class="item">😜</span>
              <span class="item">😝</span>
              <span class="item">😐</span>
              <span class="item">😶</span>
              <span class="item">😏</span>
              <span class="item">😒</span>
              <span class="item">😌</span>
              <span class="item">😔</span>
              <span class="item">😪</span>
              <span class="item">😷</span>
              <span class="item">😵</span>
              <span class="item">😎</span>
              <span class="item">😲</span>
              <span class="item">😳</span>
              <span class="item">😨</span>
              <span class="item">😰</span>
              <span class="item">😥</span>
              <span class="item">😢</span>
              <span class="item">😭</span>
              <span class="item">😱</span>
              <span class="item">😖</span>
              <span class="item">😣</span>
              <span class="item">😞</span>
              <span class="item">😓</span>
              <span class="item">😩</span>
              <span class="item">😫</span>
              <span class="item">😤</span>
              <span class="item">😡</span>
              <span class="item">😠</span>
            </div>
          </a-tab-pane>
          <a-tab-pane :key="2" title="表情包">
            <div class="emoji_list emoji_list2 scrollbar" @click="emojiClick">
              <span class="item img" v-for="item in imageList">
                <img :src="item" alt="">
              </span>
            </div>
          </a-tab-pane>
        </a-tabs>
      </div>
    </template>
  </a-trigger>
</template>

<style lang="less">
.emoji_trigger_content {
  background-color: var(--color-bg-1);
  border-radius: 5px;
  border: @f_border;
  width: 303px;

  .emoji_list{
    display: flex;
    flex-wrap: wrap;
    .item{
      width: 30px;
      height: 30px;
      display: flex;
      justify-content: center;
      align-items: center;
      font-size: 20px;
      cursor: pointer;
      &:hover{
        background-color: var(--color-fill-1);
      }
    }
  }

  .emoji_list2{
    max-height: 300px;
    overflow-y: auto;
    .item{
      width: 58px;
      height: 58px;
    }
    img{
      width: 80%;
      height: 80%;
    }
  }
}
</style>