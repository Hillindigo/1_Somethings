<script setup lang="ts">
import {bannerListApi, type bannerListType} from "@/api/banner_api";
import {reactive, ref} from "vue";
import type {listResponse} from "@/api";
import {IconClose} from "@arco-design/web-vue/es/icon";

const data = reactive<listResponse<bannerListType>>({
  list: [],
  count: 0
})

async function getData() {
  const res = await bannerListApi({type: 1})
  const list: bannerListType[] = []
  res.data.list.forEach((item) => {
    if (item.show) {
      list.push(item)
    }
  })
  data.list = list
}

function goItem(item: bannerListType) {
  if (item.href == "") {
    return
  }
  window.open(item.href, "_blank")
}

getData()
const isClose = ref(false)

function close(){
  isClose.value = true
}

</script>

<template>
  <div class="f_banner_list_com" v-if="!isClose">
    <IconClose class="close" @click="close" title="删除"></IconClose>
    <a-carousel auto-play  show-arrow="hover" class="f_banner_com" v-if="data.list.length">
      <a-carousel-item v-for="item in data.list">
        <img @click="goItem(item)" :src="item.cover" :class="{isHref: item.href !== ''}" :style="{ width: '100%' }">
      </a-carousel-item>
    </a-carousel>
  </div>

</template>

<style lang="less">
.f_banner_list_com{
  position: relative;
  &:hover{
    >.close{
      opacity: 1;
    }
  }
  >.close{
    position: absolute;
    right: 10px;
    top: 10px;
    z-index: 10;
    cursor: pointer;
    opacity: 0;
  }
  .f_banner_com {
    height: 300px;

    .arco-carousel-item-current {
      border-radius: 5px;

      img.isHref {
        cursor: pointer;
      }
    }
  }
}

</style>