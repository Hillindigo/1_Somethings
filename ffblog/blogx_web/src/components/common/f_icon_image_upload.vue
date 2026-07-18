<script setup lang="ts">
import {userStorei} from "@/stores/user_store";
import {type FileItem, Message} from "@arco-design/web-vue";
import type {baseResponse} from "@/api";
import {IconImage} from "@arco-design/web-vue/es/icon";

const store = userStorei()

const emits = defineEmits(["ok"])

function fileUploadCallback(file: FileItem){
  const res = JSON.parse(file.response) as baseResponse<string>
  if (res.code){
    Message.error(res.msg)
    return
  }
  Message.success(res.msg)
  emits("ok", res.data)
}

</script>

<template>
  <div class="f_image_upload">
   <a-upload :show-file-list="false" action="/api/images" @success="fileUploadCallback" name="file" :headers="{token: store.userInfo.token}">
      <template #upload-button>
        <IconImage></IconImage>
      </template>
    </a-upload>
  </div>
</template>

<style lang="less">
.f_image_upload {
  width: 100%;
  .arco-input-wrapper{
    display: flex;
    margin-bottom: 10px;
  }
  .arco-image{
    &.circle{
      border-radius: 50%;
    }
  }
}
</style>