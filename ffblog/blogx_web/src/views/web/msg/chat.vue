<script setup lang="ts">
import Session_list from "@/components/web/msg/session_list.vue";
import Chat_list from "@/components/web/msg/chat_list.vue";
import {useRoute} from "vue-router";
import {userStorei} from "@/stores/user_store";
import {ref} from "vue";
import type {chatSendType} from "@/api/chat_api";
import F_icon_image_upload from "@/components/common/f_icon_image_upload.vue";
import Emoji_trigger from "@/components/web/msg/emoji_trigger.vue";

const content = ref("")
const route = useRoute()
const store = userStorei()

function sendText() {
  const data: chatSendType = {
    revUserID: Number(route.query.userID),
    msgType: 1,
    msg: {
      textMsg: {
        content: content.value,
      }
    }
  }
  store.ws?.send(JSON.stringify(data))
  content.value = ""
}

function sendImage(src: string) {
  const data: chatSendType = {
    revUserID: Number(route.query.userID),
    msgType: 2,
    msg: {
      imageMsg: {
        src: src,
      }
    }
  }
  store.ws?.send(JSON.stringify(data))
  content.value = ""
}

const sessionCount = ref(1)

function getSessionCount(count: number) {
  sessionCount.value = count
}

function select(type: "emoji" | "img", val: string) {
  if (type === "emoji") {
    const t = document.querySelector(".f_textarea textarea") as HTMLTextAreaElement
    const s1 = content.value.substring(0, t.selectionStart)
    const s2 = content.value.substring(t.selectionEnd,)
    content.value = s1 + val + s2
    return
  }
  sendImage(val)
}


</script>

<template>
  <div class="chat_view">
    <template v-if="sessionCount">
      <session_list @getSessionCount="getSessionCount"></session_list>
      <div class="chat_inner" v-if="route.query.userID">
        <chat_list></chat_list>
        <div class="chat_menu">
          <div class="icons">
            <emoji_trigger @select="select"></emoji_trigger>
            <span>
           <f_icon_image_upload @ok="sendImage"></f_icon_image_upload>
          </span>
          </div>
          <div class="chat_ipt">
            <a-textarea class="f_textarea" v-model="content" @keydown.enter="sendText"
                        :auto-size="{minRows: 6, maxRows: 6}" placeholder="请输入聊天内容"></a-textarea>
            <div class="right">
              <span class="tip">按下Enter发送内容</span>
              <a-button @click="sendText" type="outline">发送</a-button>
            </div>
          </div>
        </div>
      </div>
    </template>
    <template v-else>
      <a-empty></a-empty>
    </template>

  </div>

</template>

<style lang="less">
.chat_view {
  display: flex;
  height: 100%;

  .chat_inner {
    width: calc(100% - 260px);

    .f_chat_list_com {
      height: calc(100vh - 370px);
      overflow-y: auto;
    }

    .chat_menu {
      width: 100%;
      border-top: @f_border;

      .icons {
        display: flex;

        > span {
          padding: 8px;
          color: var(--color-text-2);
          cursor: pointer;

          &:hover {
            background-color: var(--color-fill-1);
          }

          svg {
            font-size: 16px;
          }
        }
      }

      .chat_ipt {
        position: relative;

        .arco-textarea-wrapper {
          background-color: transparent;
          border: none;
        }

        .right {
          position: absolute;
          right: 10px;
          bottom: 10px;
          z-index: 2;

          .tip {
            font-size: 12px;
            color: var(--color-text-2);
            margin-right: 10px;
          }

          .arco-btn {
            border-radius: 100px;
          }
        }
      }
    }
  }
}
</style>