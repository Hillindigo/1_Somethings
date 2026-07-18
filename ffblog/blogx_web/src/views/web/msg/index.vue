<script setup lang="ts">
import F_main from "@/components/web/f_main.vue";
import F_nav from "@/components/web/f_nav.vue";
import {IconSettings} from "@arco-design/web-vue/es/icon";
import {userStorei} from "@/stores/user_store";

const store = userStorei()
store.unReadMsg()
</script>

<template>
  <div class="msg_index_view">
    <f_nav no-scroll></f_nav>
    <f_main>
      <div class="head">
        <div class="title">消息中心</div>
        <div class="settings">
          <router-link :to="{name:'msgSettings'}">
            <IconSettings></IconSettings>
            消息设置
          </router-link>
        </div>
      </div>
      <div class="menu">
        <a-badge :count="store.unMsgInfo.commentMsgCount">
          <router-link :to="{name: 'msgComment'}">评论和@</router-link>
        </a-badge>
        <a-badge :count="store.unMsgInfo.diggMsgCount">
          <router-link :to="{name: 'msgDigg'}">赞和收藏</router-link>
        </a-badge>
        <a-badge :count="store.unMsgInfo.privateMsgCount">
          <router-link :to="{name: 'msgChat'}">私信</router-link>
        </a-badge>
        <a-badge :count="store.unMsgInfo.systemMsgCount">
          <router-link :to="{name: 'msgSystem'}">系统通知</router-link>
        </a-badge>
      </div>
      <div class="body scrollbar">
        <router-view></router-view>
      </div>
    </f_main>
  </div>
</template>

<style lang="less">
.msg_index_view {
  .f_container {
    margin: 20px 0;
    background-color: var(--color-bg-1);
    border-radius: 5px;

    .head {
      padding: 20px 20px 10px 20px;
      display: flex;
      justify-content: space-between;

      .title {
        font-size: 18px;
      }

      .settings {
        color: var(--color-text-2);
        cursor: pointer;
      }
    }

    .menu {
      border-bottom: @f_border;
      padding: 10px 20px;
      >span{
        margin-right: 40px;
        .arco-badge-number{
          right: -13px;
          top: -6px;
          transform: scale(0.7);
        }
      }
      a {
        color: var(--color-text-1);

        &.router-link-active {
          color: rgb(var(--arcoblue-6));
        }
      }
    }

    .body {
      height: calc(100vh - 188px);
      overflow-y: auto;
    }
  }
}
</style>