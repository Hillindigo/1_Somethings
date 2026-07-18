<script setup lang="ts">
import Msg_base from "@/components/web/msg/msg_base.vue";
import type{ siteMsgListType} from "@/api/site_msg_api";
</script>

<template>
  <msg_base class="msg_system_view" :t="3" v-slot="{item}:{item: siteMsgListType}">
    <div class="content">
      <div class="title">{{ item.title }}</div>
      <div class="content_inner" :class="{is_read: item.isRead}">
        <a-typography-text :ellipsis="{rows: 2, css: true}">
          {{ item.content }}
        </a-typography-text>
      </div>
      <div class="link" v-if="item.linkTitle">
        <span>链接：</span> <a :href="item.linkHref">{{ item.linkTitle }}</a>
      </div>
    </div>
  </msg_base>
</template>

<style lang="less">
.msg_system_view {
  .item {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 10px 0;
    border-bottom: @f_border;

    .content {
      width: calc(100% - 160px);

      .content_inner {
        margin: 5px 0;

        .arco-typography {
          color: rgb(var(--arcoblue-6));
        }

        &.is_read {
          .arco-typography {
            color: var(--color-text-2);
          }
        }
      }

      .link {
        font-size: 12px;

        a {
          color: rgb(var(--arcoblue-4));

        }
      }
    }

  }
}
</style>