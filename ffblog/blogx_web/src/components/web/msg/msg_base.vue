<script setup lang="ts">
import {reactive} from "vue";
import type {baseResponse, listResponse} from "@/api";
import {
  siteMsgListApi,
  type siteMsgListParams,
  type siteMsgListType,
  siteMsgReadApi,
  siteMsgRemoveApi
} from "@/api/site_msg_api";
import {Message} from "@arco-design/web-vue";
import {dateFormat} from "@/utils/date";
import {IconMore} from "@arco-design/web-vue/es/icon";

interface Props {
  t: 1 | 2 | 3
}

const props = defineProps<Props>()
const data = reactive<listResponse<siteMsgListType>>({
  list: [],
  count: 0
})
const params = reactive<siteMsgListParams>({
  t: props.t,
  limit: 10,
  page: 1,
})

async function getData() {
  const res = await siteMsgListApi(params)
  if (res.code) {
    Message.error(res.msg)
    return
  }
  data.list = res.data.list
  data.count = res.data.count
}

getData()

async function handleSelect(v: string | number | Record<string, any> | undefined, id: number) {
  let res: baseResponse<string>
  const val = v as string
  if (val === "read") {
    res = await siteMsgReadApi({id: id, t: props.t})
  } else {
    res = await siteMsgRemoveApi({id: id, t: props.t})
  }
  if (res.code) {
    Message.error(res.msg)
    return
  }
  Message.success(res.msg)
  getData()
}

</script>

<template>
  <div class="f_msg_base_com">
    <div class="list">
      <div class="item" v-for="item in data.list">
        <slot :item="item"></slot>
        <div class="action">
          <span class="date">{{ dateFormat(item.createdAt) }}</span>
          <span class="more">
              <a-dropdown @select="handleSelect($event, item.id)">
           <IconMore></IconMore>
              <template #content>
                <a-doption value="read" :disabled="item.isRead">读取消息</a-doption>
                <a-doption value="delete">删除消息</a-doption>
              </template>
            </a-dropdown>
        </span>
        </div>
      </div>
      <div class="no_data" v-if="data.list.length === 0">
        <a-empty></a-empty>
      </div>
    </div>
    <div class="page" v-if="data.list.length">
      <a-pagination :page-size="params.limit" :total="data.count" show-total v-model:current="params.page"
                    @change="getData"></a-pagination>
    </div>
  </div>
</template>

<style lang="less">
.f_msg_base_com {
  .list {
    padding: 10px 20px;

    .item {
      .action {
        .date {
          margin-right: 10px;
          color: var(--color-text-2);
        }

        .more {
          cursor: pointer;
          color: var(--color-text-2);
        }
      }
    }
  }

  .page {
    display: flex;
    justify-content: center;
    padding: 20px 0;
  }
}
</style>