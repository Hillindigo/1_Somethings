<script setup lang="ts">

import F_point from "@/components/common/f_point.vue";
import router from "@/router";
import {userStorei} from "@/stores/user_store";

const store = userStorei()

function select(v: string | number | Record<string, any> | undefined) {
  const val  = v as string
  router.push({
    name: val
  })
}

if (store.isLogin) {
  store.unReadMsg()
}


</script>

<template>
  <f_point :count="store.unReadMsgAll" class="f_nav_msg_com">
    <a-dropdown class="f_nav_msg_dropdown" @select="select">
      <a href="javascript:void 0">消息</a>
      <template #content>
        <a-doption value="msgComment">
          <f_point :count="store.unMsgInfo.commentMsgCount">评论和@</f_point>
        </a-doption>
        <a-doption value="msgDigg">
          <f_point :count="store.unMsgInfo.diggMsgCount">赞和收藏</f_point>
        </a-doption>
        <a-doption value="msgChat">
          <f_point :count="store.unMsgInfo.privateMsgCount">用户私信</f_point>
        </a-doption>
        <a-doption value="msgSystem">
          <f_point :count="store.unMsgInfo.systemMsgCount">系统通知</f_point>
        </a-doption>
      </template>
    </a-dropdown>
  </f_point>
</template>

<style lang="less">
.f_nav_msg_dropdown {
  .arco-dropdown-option-content {
    padding: 0 6px;
  }
}
</style>