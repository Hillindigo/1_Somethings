<script setup lang="ts">
import {MdPreview} from "md-editor-v3";
import type {msgType} from "@/api/chat_api";
import {computed} from "vue";

interface Props {
  msg: msgType
}

const props = defineProps<Props>()
const msgType = computed(()=>{
  if (props.msg.textMsg){
    return 1
  }
  if (props.msg.imageMsg){
    return 2
  }
  if (props.msg.markdownMsg){
    return 3
  }
  return 0
})


</script>

<template>
  <div class="f_msg_preview_com" :class="`msg_${msgType}`">
    <template v-if="msgType === 1">
      {{ props.msg?.textMsg?.content }}
    </template>
    <template v-else-if="msgType === 2">
      [图片消息]
    </template>
    <template v-else-if="msgType === 3">
     [图文消息]
    </template>
    <template v-else>
      [未知消息]
    </template>
  </div>
</template>

<style  lang="less">

</style>