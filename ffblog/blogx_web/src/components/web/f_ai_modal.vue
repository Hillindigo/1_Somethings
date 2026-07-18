<script setup lang="ts">
import {nextTick, reactive, ref} from "vue";
import {aiSiteInfoApi, type aiSiteInfoType} from "@/api/site_api";
import {Message} from "@arco-design/web-vue";
import {userStorei} from "@/stores/user_store";
import {showLogin} from "@/components/web/f_login";
import type {baseResponse} from "@/api";

const store = userStorei()

interface Props {
  visible: boolean
}

interface chatResponse {
  isMe: boolean
  nickname: string
  avatar: string
  content: string
}

const chatList = ref<chatResponse[]>([])

const props = defineProps<Props>()
const emits = defineEmits(["update:visible"])

function cancel() {
  chatList.value = []
  emits("update:visible", false)
}


async function send() {
  if (!store.isLogin) {
    Message.warning("请登录")
    showLogin()
    return
  }
  chatList.value.push({
    isMe: true,
    nickname: store.userInfo.nickName,
    avatar: store.userInfo.avatar,
    content: key.value,
  })
  const item = reactive<chatResponse>({
    isMe: false,
    nickname: aiData.nickname,
    avatar:aiData.avatar,
    content: "",
  })

  const eventSource = new EventSource(`/api/ai/article?content=${key.value}&token=${store.userInfo.token}`)
  eventSource.onmessage = (e) => {
    const message = JSON.parse(e.data) as baseResponse<string>;
    item.content += message.data
    console.log(message)
  };
  eventSource.onerror = (e) => {
    console.log(e)
    eventSource.close()
    return
  };
  //处理服务器响应报文中的自定义事件
  eventSource.addEventListener("close", function (e) {
    console.log(e)
  });

  chatList.value.push(item)

  key.value = ""

}

const textareaRef = ref()
const aiData = reactive<chatResponse>({
  isMe: false,
  nickname: "",
  avatar: "",
  content: "",
})

async function beforeOpen() {
  chatList.value = []
  aiData.content = ""
  const res = await aiSiteInfoApi()
  if (res.code) {
    Message.error(res.msg)
    return
  }
  Object.assign(aiData, res.data)
  chatList.value.push(aiData)
  const list = res.data.abstract.split("")
  for (let i = 0; i < list.length; i++) {
    setTimeout(() => {
      aiData.content += list[i]
    }, 40 * i)
  }

  nextTick(() => {
    textareaRef.value.focus()
  })
}

const key = ref("")

</script>

<template>
  <a-modal @before-open="beforeOpen" :footer="false" modal-class="f_ai_modal"
           :visible="props.visible" @cancel="cancel">
    <div class="body scrollbar">
      <div class="item" :class="{isMe: item.isMe}" v-for="item in chatList">
        <a-avatar :image-url="item.avatar"></a-avatar>
        <div class="chat" v-if="!item.isMe" v-html="item.content"></div>
        <div class="chat" v-else>{{ item.content }}</div>
      </div>
    </div>
    <div class="menu">
      <a-textarea ref="textareaRef" :auto-size="{minRows: 4, maxRows: 4}" v-model="key" @keydown.enter="send"
                  placeholder="请输入你感兴趣的内容"></a-textarea>
      <div class="info">
        <span>按Enter发送，按Ctrl+Enter换行</span>
        <a-button @click="send" type="primary">发送</a-button>
      </div>
    </div>
  </a-modal>
</template>

<style lang="less">
.f_ai_modal {
  .arco-modal-header {
    display: none;
  }

  .arco-modal-body {
    padding: 0;
  }

  .body {
    min-height: 40vh;
    max-height: 60vh;
    overflow-y: auto;
    padding: 20px;

    .item {
      display: flex;
      margin-bottom: 10px;

      .arco-avatar {
        flex-shrink: 0;
        margin-right: 10px;
      }

      .chat {
        background-color: var(--color-fill-2);
        padding: 10px;
        border-radius: 5px;
        color: var(--color-text-1);
        position: relative;

        &::after {
          position: absolute;
          left: -15px;
          top: 10px;
          width: 0;
          height: 0;
          border-width: 8px;
          border-style: solid;
          border-color: transparent;
          border-right-color: var(--color-fill-2);
          content: "";
          display: block;
        }
      }

      &.isMe {
        justify-content: end;
        flex-direction: row-reverse;

        .arco-avatar {
          margin-right: 0;
          margin-left: 10px;
        }

        .chat {
          &::after {
            position: absolute;
            left: inherit;
            right: -15px;
            top: 10px;
            width: 0;
            height: 0;
            border-width: 8px;
            border-style: solid;
            border-color: transparent;
            border-left-color: var(--color-fill-2);
            content: "";
            display: block;
          }
        }
      }
    }
  }

  .menu {
    padding: 10px 20px;
    border-top: @f_border;
    position: relative;

    .info {
      position: absolute;
      right: 30px;
      bottom: 25px;
      z-index: 2;

      > span {
        font-size: 12px;
        margin-right: 10px;
        color: var(--color-text-2);
      }
    }
  }
}
</style>