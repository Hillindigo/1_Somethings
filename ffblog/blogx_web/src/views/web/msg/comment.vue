<script setup lang="ts">
import Msg_base from "@/components/web/msg/msg_base.vue";
import type{ siteMsgListType} from "@/api/site_msg_api";
import F_label from "@/components/common/f_label.vue";
import {relationOptions, siteMsgTypeOptions} from "@/options/options";
import {goUser} from "@/utils/go_router";

</script>

<template>
  <msg_base class="msg_comment_view" :t="1" v-slot="{item}:{item: siteMsgListType}">
      <a-avatar @click="goUser(item.ActionUserID)" :image-url="item.actionUserAvatar"></a-avatar>
      <div class="content">
        <div class="user">
          <span  @click="goUser(item.ActionUserID)" class="nick">{{ item.actionUserNickname }}</span>
          <span class="relation">
               <f_label v-if="item.relation !== 0 && item.relation !== 1" :options="relationOptions"
                        :value="item.relation"></f_label>
          </span>
          <span class="type">
            <f_label type="span" :options="siteMsgTypeOptions" :value="item.type"></f_label>
          </span>
        </div>
        <div class="content_inner" :class="{is_read: item.isRead}">
          <a-typography-text :ellipsis="{rows: 2, css: true}">
            {{ item.content }}
          </a-typography-text>
        </div>
        <div class="article">
          <router-link :to="{name: 'articleDetail', params: {id: item.articleID}}">{{ item.articleTitle }}</router-link>
        </div>
      </div>
  </msg_base>
</template>

<style lang="less">
.msg_comment_view {
  .item {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 10px 0;
    border-bottom: @f_border;

    .arco-avatar {
      flex-shrink: 0;
    }

    .content {
      width: calc(100% - 160px);

      .user {
        .relation {
          margin: 0 10px;
        }

        .type {
          color: var(--color-text-2);
          font-size: 12px;
        }
      }

      .content_inner {
        margin: 5px 0;
        .arco-typography{
          color: rgb(var(--arcoblue-6));
        }

        &.is_read{
          .arco-typography{
            color: var(--color-text-2);
          }
        }
      }

      .article {
        a {
          color: var(--color-text-2);
          font-size: 12px;
        }
      }
    }

  }
}
</style>