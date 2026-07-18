<script setup lang="ts">
import {reactive} from "vue";
import {siteMsgSettingsApi, type siteMsgSettingsType, siteMsgSettingsUpdateApi} from "@/api/site_msg_api";
import {Message} from "@arco-design/web-vue";

const data = reactive<siteMsgSettingsType>({
  openCommentMessage: false,
  openDiggMessage: false,
  openPrivateChat: false,
})

async function getData() {
  const res = await siteMsgSettingsApi()
  if (res.code) {
    Message.error(res.msg)
    return
  }
  Object.assign(data, res.data)
}

getData()

async function update() {
  const res = await siteMsgSettingsUpdateApi(data)
  if (res.code) {
    Message.error(res.msg)
    return
  }
  Message.success(res.msg)
  getData()
}


</script>

<template>
  <div class="msg_settings_view">
    <a-form  :label-col-props="{span: 2}">
      <a-form-item label="评论和@">
        <a-switch @change="update" v-model="data.openCommentMessage"></a-switch>
        <template #help>有人对我的内容进行评论时，我会收到消息提醒</template>
      </a-form-item>
      <a-form-item label="赞和收藏">
        <a-switch @change="update" v-model="data.openDiggMessage"></a-switch>
        <template #help>有人对我表达的内容表示赞同时，我会收到消息提醒</template>
      </a-form-item>
      <a-form-item label="陌生人私信">
        <a-switch @change="update" v-model="data.openPrivateChat"></a-switch>
        <template #help>开启后陌生人可以私信我</template>
      </a-form-item>
    </a-form>
  </div>

</template>

<style lang="less">
.msg_settings_view{
  padding: 10px 20px;
}
</style>