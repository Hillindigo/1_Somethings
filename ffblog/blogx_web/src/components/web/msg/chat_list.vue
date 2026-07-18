<script setup lang="ts">
import {nextTick, reactive, ref, watch} from "vue";
import type {listResponse} from "@/api";
import {chatListApi, type chatListParams, type chatListType} from "@/api/chat_api";
import {Message} from "@arco-design/web-vue";
import {useRoute} from "vue-router";
import {dateTimeFormat} from "../../../utils/date";
import {userStorei} from "@/stores/user_store";
import Msg from "@/components/web/msg/msg.vue";
import {chatRemoveApi} from "@/api/chat_api";
import {goUser} from "@/utils/go_router";

const store = userStorei()
const route = useRoute()
const chat = reactive<listResponse<chatListType>>({
  list: [],
  count: 0
})
const params = reactive<chatListParams>({
  revUserID: 0,
  type: 1,
  page: 1,
})

async function getData() {
  params.revUserID = Number(route.query.userID)
  const res = await chatListApi(params)
  if (res.code) {
    Message.error(res.msg)
    return
  }
  chat.list = res.data.list.reverse()
  chat.count = res.data.count
  nextTick(goBottom)
}

watch(() => route.query.userID, () => {
  getData()
}, {immediate: true})

function goBottom() {
  const ele = document.querySelector(".f_chat_list_com ") as HTMLDivElement
  const top = ele.scrollHeight
  ele.scrollTo({
    top: top,
    behavior: "smooth"
  })
}

watch(() => store.wsChatList, () => {
  const data = store.wsChatList[0]
  chat.list.push(data)
  chat.count++
  nextTick(goBottom)
})

async function more() {
  (params.page as number)++
  const res = await chatListApi(params)
  if (res.code) {
    Message.error(res.msg)
    return
  }
  if (res.data.list.length === 0) {
    Message.warning("已经是最早的消息啦")
    return
  }
  chat.list = [...res.data.list.reverse(), ...chat.list]
  chat.count += res.data.count
}

const checkIDList = ref<number[]>([])
const isCheck = ref(false)

function check() {
  isCheck.value = !isCheck.value
}

function inList(id: number): boolean {
  const item = checkIDList.value.find((item) => item === id)
  if (item) {
    return true
  }
  return false
}

async function remove(){
  const res = await chatRemoveApi(checkIDList.value)
  if (res.code){
    Message.error(res.msg)
    return
  }
  Message.success(res.msg)
  checkIDList.value = []
  getData()
}

</script>

<template>
  <div class="f_chat_list_com scrollbar">
    <div class="actions">
      <span class="plcz" @click="check">批量操作</span>
      <a-button status="danger" @click="remove" v-if="isCheck && checkIDList.length" size="mini">批量删除</a-button>
    </div>
    <div class="inner">
      <div class="more">
        <span @click="more">加载更多</span>
      </div>
      <a-checkbox-group v-model="checkIDList">
        <div class="item" :class="{isMe: item.isMe, isCheck: inList(item.id)}" v-for="item in chat.list">
          <div class="top">
            <div class="date">{{ dateTimeFormat(item.createdAt) }}</div>
          </div>
          <div class="bottom">
            <a-checkbox v-if="isCheck" :value="item.id"></a-checkbox>
            <a-avatar @click="goUser( item.sendUserID)" :image-url="item.sendUserAvatar"></a-avatar>
            <div class="content">
              <msg :msg="item.msg"></msg>
            </div>
          </div>
        </div>
      </a-checkbox-group>
    </div>
  </div>
</template>

<style lang="less">

.f_chat_list_com {
  .actions {
    display: flex;
    align-items: center;
    height: 35px;
    padding: 10px 20px;

    .plcz {
      cursor: pointer;
      font-size: 12px;
      color: var(--color-text-2);
      margin-right: 10px;
    }
  }

  .more {
    display: flex;
    justify-content: center;
    font-size: 12px;
    color: var(--color-text-2);

    span {
      cursor: pointer;
    }
  }

  .arco-checkbox-group {
    width: 100%;
  }

  .item {
    padding: 10px 20px;

    &.isMe {
      .bottom {
        flex-direction: row-reverse;

        .content {
          margin-left: 0;
          margin-right: 10px;

          &::after {
            left: inherit;
            right: -15px;
            border-color: transparent;
            border-left-color: var(--color-fill-2);
          }
        }
      }

    }

    &.isCheck{
      background-color: var(--color-fill-1);
    }

    .top {
      display: flex;
      justify-content: center;
      font-size: 12px;
      color: var(--color-text-2);
    }

    .bottom {
      margin-top: 5px;
      display: flex;

      .arco-avatar {
        flex-shrink: 0;
      }

      .content {
        margin-left: 10px;
        background-color: var(--color-fill-2);
        padding: 10px;
        border-radius: 5px;
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
    }
  }
}
</style>