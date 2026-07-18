<script setup lang="ts">

import {relationOptions} from "@/options/options";
import {dateFormat} from "@/utils/date";
import {reactive} from "vue";
import type {listResponse} from "@/api";
import {sessionListApi, type sessionListType, sessionRemoveApi} from "@/api/chat_api";
import {Message} from "@arco-design/web-vue";
import F_label from "@/components/common/f_label.vue";
import router from "@/router";
import {useRoute} from "vue-router";
import {IconImage} from "@arco-design/web-vue/es/icon";
import {userBaseStorei} from "@/stores/user_base_store";
import Msg_preview from "@/components/web/msg/msg_preview.vue";
import {goUser} from "@/utils/go_router";

const baseStore = userBaseStorei()
const route = useRoute()
const session = reactive<listResponse<sessionListType>>({
  list: [],
  count: 0
})
const emits = defineEmits(["getSessionCount"])

async function getData() {
  const res = await sessionListApi({})
  if (res.code) {
    Message.error(res.msg)
    return
  }
  Object.assign(session, res.data)

  // 判断route上的用户，在不在列表里面
  const userID = Number(route.query.userID)
  if (!isNaN(userID)) {
    initUser(userID)
  }
  emits("getSessionCount", res.data.list.length)

}

function initUser(userID: number) {
  const item = session.list.find((item) => item.userID === userID)
  if (item) {
    return
  }
  if (baseStore.userBase.userID) {
    session.list.unshift({
      userID: userID,
      userNickname: baseStore.userBase.nickname,
      userAvatar: baseStore.userBase.avatar,
      newMsgDate: dateFormat(new Date().toString()),
      relation: baseStore.userBase.relation,
      msgType: 1,
      msg: {
        textMsg: {
          content: "",
        }
      }
    })
  }

}

getData()

function goItem(item: sessionListType) {
  router.push({
    name: "msgChat",
    query: {
      userID: item.userID,
    }
  })
}


async function removeSession(item: sessionListType) {
  const res = await sessionRemoveApi(item.userID)
  if (res.code) {
    Message.error(res.msg)
    return
  }
  Message.success(res.msg)
  getData()

  const div = document.querySelector(`.contextMenuTrigger_${item.userID}`) as HTMLDivElement
  if (div){
    div.remove()
  }
}

</script>

<template>
  <div class="f_session_list_com">
    <div class="item" :class="{active: Number(route.query.userID) === item.userID}" @click="goItem(item)"
         v-for="item in session.list">
      <a-trigger :class="`contextMenuTrigger_${item.userID}`" trigger="contextMenu" align-point>
        <a-avatar @click.stop="goUser(item.userID)" :image-url="item.userAvatar"></a-avatar>
        <template #content>
          <div class="item_context_menu_user">
            <span @click="removeSession(item)">删除会话</span>
          </div>
        </template>
      </a-trigger>

      <div class="info">
        <div class="top">
          <div class="left">
            <a-typography-text :ellipsis="{rows: 1}">{{ item.userNickname }}</a-typography-text>
            <f_label :options="relationOptions" :value="item.relation"></f_label>
          </div>
          <div class="date">{{ dateFormat(item.newMsgDate) }}</div>
        </div>
        <div class="bottom">
          <a-typography-text :ellipsis="{rows: 1}">
            <msg_preview :msg="item.msg"></msg_preview>
          </a-typography-text>
        </div>
      </div>
    </div>
  </div>
</template>

<style lang="less">
.item_context_menu_user {
  background-color: var(--color-bg-1);
  padding: 20px 0;

  span {
    color: var(--color-text-2);
    cursor: pointer;
    padding: 10px 20px;

    &:hover {
      background-color: var(--color-fill-1);
    }
  }
}

.f_session_list_com {
  width: 260px;
  border-right: @f_border;
  height: 100%;

  .item {
    display: flex;
    padding: 10px 20px;

    .arco-avatar {
      flex-shrink: 0;
    }

    .arco-tag-size-medium {
      transform: scale(0.8);
    }

    &:hover {
      background-color: var(--color-fill-1);
    }

    &.active {
      background-color: var(--color-fill-1);
    }

    .info {
      width: 100%;
      margin-left: 10px;
      display: flex;
      flex-direction: column;
      justify-content: start;

      .top {
        display: flex;
        justify-content: space-between;
        align-items: center;


        .left {
          display: flex;
          align-items: center;
        }

        .date {
          font-size: 12px;
          color: var(--color-text-2);
        }
      }
    }
  }
}

</style>